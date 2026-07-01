#!/usr/bin/env python
"""Entry point compatibile con la vecchia invocazione:

    python generate_page.py playnite.csv -o index.html --title "PLAYNITE" --config config.ini

La logica vera e propria vive nel pacchetto playnite_page/.
"""
from playnite_page.cli import main

if __name__ == "__main__":
    raise SystemExit(main())
