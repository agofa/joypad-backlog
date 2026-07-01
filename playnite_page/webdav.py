"""Upload opzionale del file generato su un server WebDAV.

Variabili d'ambiente attese (in un file .env nella cartella di lavoro
o già presenti nell'ambiente):

    WEBDAV_URL       URL del server WebDAV (obbligatoria)
    WEBDAV_USERNAME  utente (obbligatoria)
    WEBDAV_PASSWORD  password (obbligatoria)
    WEBDAV_PATH      percorso remoto del file (opzionale, default: nome del file)
"""
from __future__ import annotations

import os
from pathlib import Path


class WebDAVConfigError(RuntimeError):
    """Sollevata quando mancano variabili d'ambiente necessarie per l'upload."""


def upload_webdav(file_path: Path) -> None:
    from dotenv import load_dotenv
    from webdav3.client import Client

    load_dotenv()

    hostname = os.getenv("WEBDAV_URL")
    username = os.getenv("WEBDAV_USERNAME")
    password = os.getenv("WEBDAV_PASSWORD")

    missing = [name for name, val in (
        ("WEBDAV_URL", hostname),
        ("WEBDAV_USERNAME", username),
        ("WEBDAV_PASSWORD", password),
    ) if not val]
    if missing:
        raise WebDAVConfigError(
            "Variabili d'ambiente mancanti per l'upload WebDAV: " + ", ".join(missing)
        )

    client = Client({
        "webdav_hostname": hostname,
        "webdav_login": username,
        "webdav_password": password,
    })
    client.verify = True  # Usa False solo con certificati SSL autofirmati

    remote_path = os.getenv("WEBDAV_PATH", "").lstrip("/") or file_path.name

    try:
        client.upload_sync(remote_path=remote_path, local_path=str(file_path))
        print(f"✅ Upload completato su WebDAV: {remote_path}")
    except Exception as exc:  # noqa: BLE001 - vogliamo un messaggio comprensibile all'utente
        print(f"❌ Errore upload WebDAV: {exc}")
