"""Interfaccia a riga di comando del generatore di pagine."""
from __future__ import annotations

import argparse
import sys
from datetime import datetime
from pathlib import Path
from typing import List, Optional

from .config import load_config
from .parsers import autodetect_and_parse
from .render import render_html
from .utils import format_seconds, parse_datetime, seconds_from_raw


def build_arg_parser() -> argparse.ArgumentParser:
    ap = argparse.ArgumentParser(
        prog="generate_page.py",
        description="Genera una pagina HTML (tabella + schede) a partire da un export della libreria giochi.",
    )
    ap.add_argument("input", type=Path, help="File di input (CSV, JSON, INI, o formato custom a blocchi)")
    ap.add_argument("-o", "--output", type=Path, default=Path("output.html"), help="File HTML di output")
    ap.add_argument("--title", default="Pagina generata", help="Titolo della pagina HTML")
    ap.add_argument("--template", type=Path, help="Template HTML personalizzato (opzionale)")
    ap.add_argument("--config", type=Path, help="File config.ini con ordine/colonne nascoste (opzionale)")
    ap.add_argument("--upload", action="store_true", help="Carica il file generato su WebDAV (richiede variabili d'ambiente, vedi README)")
    return ap


def _normalize_rows(rows: list[dict], config: dict) -> tuple[int, int]:
    """Applica formattazione date e Time Played in-place. Ritorna (totale_secondi, n_righe)."""
    for row in rows:
        for col in config.get("dates", []):
            if col in row and row[col]:
                dt = parse_datetime(row[col])
                if dt:
                    row[col] = dt.strftime("%d/%m/%Y")
                    row[col + "__parsed"] = dt

    total_seconds = 0
    for row in rows:
        seconds = seconds_from_raw(row.get("Time Played", ""))
        total_seconds += seconds
        row["Time Played__seconds"] = seconds
        row["Time Played"] = format_seconds(seconds)

    return total_seconds, len(rows)


def main(argv: Optional[List[str]] = None) -> int:
    args = build_arg_parser().parse_args(argv)

    if not args.input.exists():
        print(f"Errore: file non trovato {args.input}", file=sys.stderr)
        return 2

    config = load_config(args.config)
    source_type, rows, headers = autodetect_and_parse(args.input)

    total_seconds, _ = _normalize_rows(rows, config)
    config["total_played"] = format_seconds(total_seconds)
    config["total_days"] = str(total_seconds // 86400)

    if rows and "Name" in rows[0]:
        rows.sort(key=lambda r: r.get("Name", "").lower())

    today = datetime.now().strftime("%d/%m/%Y %H:%M")
    html = render_html(
        title=args.title,
        input_name=args.input.name,
        source_type=source_type,
        rows=rows,
        headers=headers,
        template_path=args.template,
        config=config,
        today=today,
    )
    args.output.write_text(html, encoding="utf-8")
    print(f"✅ Generato: {args.output.resolve()}")

    if args.upload:
        from .webdav import upload_webdav
        upload_webdav(args.output)

    return 0
