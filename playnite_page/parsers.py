"""Parser per i diversi formati di file di input.

Ogni parser restituisce una tupla ``(source_type, rows, headers)`` dove:
  - ``source_type``: etichetta del formato ("csv", "json", "ini", "custom")
  - ``rows``: lista di dizionari, una entry per gioco/record
  - ``headers``: lista ordinata delle colonne disponibili
"""
from __future__ import annotations

import json
import re
import csv
from pathlib import Path
from configparser import ConfigParser
from typing import Any, Dict, List, Tuple

ParseResult = Tuple[str, List[Dict[str, Any]], List[str]]


def parse_json(path: Path) -> ParseResult:
    data = json.loads(path.read_text(encoding="utf-8"))
    if isinstance(data, dict):
        rows = [data]
    elif isinstance(data, list):
        rows = data if (data and isinstance(data[0], dict)) else [{"value": v} for v in data]
    else:
        rows = [{"value": data}]
    return "json", rows, list(rows[0].keys()) if rows else []


def parse_csv(path: Path) -> ParseResult:
    with path.open(newline="", encoding="utf-8-sig") as f:
        reader = csv.DictReader(f)
        rows = [{k.strip(): v for k, v in r.items()} for r in reader]
        headers = [h.strip() for h in (reader.fieldnames or [])]
    return "csv", rows, headers


def parse_ini(path: Path) -> ParseResult:
    cp = ConfigParser()
    cp.read(path, encoding="utf-8")
    rows = []
    for section in cp.sections():
        row = {"__section__": section}
        row.update(cp.items(section))
        rows.append(row)
    return "ini", rows, list(rows[0].keys()) if rows else []


_CUSTOM_BLOCK_SEP = re.compile(r"^---\s*$", re.M)
_KEY_VAL = re.compile(r"^\s*([^:#\n]+?)\s*:\s*(.*?)\s*$")


def parse_custom_blocks(path: Path) -> ParseResult:
    """Formato di fallback: blocchi separati da '---' con righe 'chiave: valore'."""
    text = path.read_text(encoding="utf-8").strip()
    parts = [s.strip() for s in re.split(_CUSTOM_BLOCK_SEP, text) if s.strip()]
    rows: List[Dict[str, Any]] = []
    for part in parts:
        current_key = None
        row: Dict[str, Any] = {}
        for raw_line in part.splitlines():
            m = _KEY_VAL.match(raw_line)
            if m:
                current_key = m.group(1).strip()
                row[current_key] = m.group(2).strip()
            elif current_key:
                row[current_key] += "\n" + raw_line.rstrip()
            else:
                row["testo"] = row.get("testo", "") + ("\n" if row.get("testo") else "") + raw_line.rstrip()
        if row:
            rows.append(row)
    return "custom", rows, list(rows[0].keys()) if rows else []


def autodetect_and_parse(path: Path) -> ParseResult:
    """Sceglie il parser in base all'estensione del file."""
    ext = path.suffix.lower()
    if ext == ".json":
        return parse_json(path)
    if ext == ".csv":
        return parse_csv(path)
    if ext in {".ini", ".cfg"}:
        return parse_ini(path)
    return parse_custom_blocks(path)
