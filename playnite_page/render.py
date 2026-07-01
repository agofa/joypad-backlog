"""Rendering della pagina HTML tramite Jinja2."""
from __future__ import annotations

from pathlib import Path
from typing import Any, Dict, List, Optional

from jinja2 import Template

from .config import Config

# Template minimale usato quando non viene passato --template.
# Per un layout completo con card, ricerca e ordinamento vedi template.html
# nella root del progetto.
DEFAULT_TEMPLATE_PATH = Path(__file__).parent / "default_template.html"


def _load_default_template() -> str:
    return DEFAULT_TEMPLATE_PATH.read_text(encoding="utf-8")


def build_columns(headers: List[str], config: Config) -> List[str]:
    """Determina l'ordine finale delle colonne applicando order/hide da config.

    'Name' è sempre esclusa dall'elenco colonne perché viene renderizzata
    separatamente come prima colonna della tabella.
    """
    hidden = set(config.get("hide", [])) | {"Name"}
    preferred = config.get("order", [])
    columns = [c for c in preferred if c in headers and c not in hidden]
    columns += [h for h in headers if h not in columns and h not in hidden]
    return columns


def render_html(
    title: str,
    input_name: str,
    source_type: str,
    rows: List[Dict[str, Any]],
    headers: List[str],
    template_path: Optional[Path],
    config: Config,
    today: str,
) -> str:
    template_text = template_path.read_text(encoding="utf-8") if template_path else _load_default_template()
    columns = build_columns(headers, config)

    return Template(template_text).render(
        title=title,
        input_name=input_name,
        source_type=source_type,
        rows=rows,
        columns=columns,
        table=bool(columns),
        cards=True,
        total_played=config.get("total_played"),
        total_days=config.get("total_days"),
        today=today,
    )
