#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")"

# if [[ ! -d "game" ]]; then
#   python -m venv game
# fi
# source game/bin/activate
# pip install -q -r requirements.txt

INPUT="playnite.csv"

if [[ ! -f "$INPUT" ]]; then
  echo "Errore: playnite.csv non trovato in $INPUT!"
  exit 1
fi
if [[ ! -f "config.ini" ]]; then
  echo "Errore: config.ini non trovato!"
  exit 1
fi
if [[ ! -f "generate_page.py" ]]; then
  echo "Errore: generate_page.py non trovato!"
  exit 1
fi

# Aggiungere --upload per caricare automaticamente su WebDAV (vedi README).
python generate_page.py "$INPUT" -o index.html --title "PLAYNITE" --config config.ini
