"""Caricamento del file di configurazione opzionale (config.ini).

Formato atteso::

    [columns]
    order = Name,Categories,Added,Release Date,...
    hide = Name,Genres,Categories,Description
    dates = Added,Release Date,Last Played
"""
from __future__ import annotations

from configparser import ConfigParser
from pathlib import Path
from typing import Dict, List, Optional, TypedDict


class Config(TypedDict, total=False):
    order: List[str]
    hide: List[str]
    dates: List[str]
    total_played: str
    total_days: str


def _split_csv_option(cp: ConfigParser, section: str, option: str) -> List[str]:
    if not cp.has_option(section, option):
        return []
    return [c.strip() for c in cp.get(section, option).split(",") if c.strip()]


def load_config(path: Optional[Path]) -> Config:
    """Legge order/hide/dates dalla sezione [columns]. Ritorna liste vuote se assenti."""
    empty: Config = {"order": [], "hide": [], "dates": []}
    if not path or not path.exists():
        return empty

    cp = ConfigParser()
    cp.read(path, encoding="utf-8")
    if not cp.has_section("columns"):
        return empty

    return {
        "order": _split_csv_option(cp, "columns", "order"),
        "hide": _split_csv_option(cp, "columns", "hide"),
        "dates": _split_csv_option(cp, "columns", "dates"),
    }
