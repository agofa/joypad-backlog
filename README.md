# 🕹️ Generatore di Pagina HTML per Videogiochi

Genera una pagina HTML statica (tabella ordinabile/filtrabile + schede) a
partire da un export della libreria giochi di **Playnite** (CSV), con
supporto anche a JSON, INI e a un formato custom a blocchi. Può caricare
automaticamente il risultato su un server WebDAV.

## 🎮 Prerequisito Playnite

Installa l'estensione **"Library Exporter Advanced"** di darklinkpower.
Permette di esportare la libreria Playnite in CSV con opzione per
configurare quali campi includere.

## 📁 Struttura del progetto

```
.
├── generate_page.py          # entry point CLI (compatibile con lista.sh)
├── playnite_page/            # logica applicativa
│   ├── parsers.py            # lettura CSV / JSON / INI / formato custom
│   ├── config.py             # lettura di config.ini (ordine/colonne nascoste)
│   ├── utils.py              # parsing date e formattazione Time Played
│   ├── render.py             # rendering Jinja2 della pagina HTML
│   ├── webdav.py             # upload opzionale su WebDAV
│   └── default_template.html # template HTML usato se non se ne passa uno custom
├── template.html.example      # template alternativo più curato
├── config.ini                 # esempio di configurazione colonne
├── lista.sh                   # script di lancio (venv + generazione)
└── requirements.txt
```

## ⚙️ Installazione

Di seguito un esempio, da adeguare al proprio sistema:

```bash
python -m venv joypad
source joypad/bin/activate      # su Windows: joypad\Scripts\activate
pip install -r requirements.txt
```

Se si aggiungono o aggiornano dipendenze, rigenerare il file dei requisiti con:

```bash
pip freeze > requirements.txt
```

## 🚀 Utilizzo

```bash
python generate_page.py <file_input> [-o output.html] [--title "Titolo"] [--template template.html] [--config config.ini] [--upload]
```

| Parametro | Descrizione |
|---|---|
| `file_input` | File di dati in ingresso (obbligatorio): `.csv`, `.json`, `.ini`/`.cfg`, oppure formato custom a blocchi. |
| `-o`, `--output` | Percorso del file HTML generato (default: `output.html`). |
| `--title` | Titolo della pagina HTML. |
| `--template` | Template HTML personalizzato (Jinja2). Se omesso, usa il template incluso in `playnite_page/default_template.html`. |
| `--config` | File `config.ini` per ordinare/nascondere colonne (vedi sotto). |
| `--upload` | Dopo la generazione, carica il file su WebDAV (richiede le variabili d'ambiente descritte più sotto). |

### Esempio con lo script Playnite

```bash
./lista.sh
```

`lista.sh` verifica la presenza dei file necessari (`playnite.csv`,
`config.ini`, `generate_page.py`) nella cartella corrente e lancia la
generazione. Attiva prima il virtualenv (vedi Installazione) se non lo hai
già fatto in questa sessione della shell.

## 🗂️ config.ini — ordinare e nascondere colonne

```ini
[columns]
order = Name,Categories,Added,Release Date,Last Played,Time Played,Completion Status,Platforms,Sources,Community Score,Critic Score
hide = Name,Genres,Categories,Description
dates = Added,Release Date,Last Played
```

- `order`: ordine con cui compaiono le colonne nella tabella (le colonne non elencate vengono aggiunte in coda).
- `hide`: colonne da escludere dalla tabella (`Name` è sempre esclusa dall'elenco colonne perché viene mostrata come prima colonna).
- `dates`: colonne da interpretare come date (`gg/mm/aaaa` o `gg/mm/aaaa hh:mm:ss`) e riformattare in `gg/mm/aaaa`.

## ☁️ Configurazione WebDAV (opzionale)

Per abilitare l'upload automatico (`--upload`), creare un file `.env` nella
directory principale del progetto con queste variabili:

```
WEBDAV_URL=https://tuo.server/webdav/path/
WEBDAV_USERNAME=tuo_user
WEBDAV_PASSWORD=tua_password
WEBDAV_PATH=cartella/remota/index.html   # opzionale, default: nome del file locale
```

> ⚠️ Non versionare il file `.env` (è già in `.gitignore`): contiene credenziali.

## 💡 Suggerimenti

- Si può personalizzare completamente il layout HTML scrivendo un template Jinja2 e passandolo con `--template` (vedi `template.html.example` come punto di partenza).
- La colonna "Time Played" viene normalizzata da secondi grezzi a formato `h:mm` e resta ordinabile numericamente (tramite l'attributo `data-seconds` nella cella).
- La colonna "#" è un contatore ricalcolato automaticamente dopo ogni ordinamento e non è ordinabile.
- `playnite.csv` non va versionato (è in `.gitignore`): contiene la tua libreria giochi personale.
