"""Piccole utility di formattazione/parsing usate dal generatore."""
from __future__ import annotations

from datetime import datetime
from typing import Optional

_DATE_FORMATS = ("%d/%m/%Y %H:%M:%S", "%d/%m/%Y")


def parse_datetime(value: str) -> Optional[datetime]:
    """Prova a interpretare 'gg/mm/aaaa' o 'gg/mm/aaaa hh:mm:ss'."""
    for fmt in _DATE_FORMATS:
        try:
            return datetime.strptime(value, fmt)
        except ValueError:
            continue
    return None


def format_seconds(seconds: int) -> str:
    """Converte secondi totali in stringa 'h:mm' (anche oltre le 24h)."""
    hours, remainder = divmod(seconds, 3600)
    minutes, _ = divmod(remainder, 60)
    return f"{hours}:{minutes:02d}"


def seconds_from_raw(raw: str) -> int:
    """Converte il campo 'Time Played' (stringa numerica di secondi) in int, tollerante a valori vuoti/non validi."""
    raw = (raw or "").strip()
    if not raw:
        return 0
    try:
        return int(float(raw))
    except ValueError:
        return 0
