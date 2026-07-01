# Changelog

Tutte le modifiche rilevanti al progetto sono documentate in questo file.

Il formato è basato su [Keep a Changelog](https://keepachangelog.com/it/1.0.0/),
e il progetto segue [Semantic Versioning](https://semver.org/lang/it/).

## [1.0.0] - 2026-07-01

### Added
- Prima release pubblica del generatore di pagine HTML per librerie di videogiochi Playnite.
- Parser per file di input CSV, JSON, INI e formato custom a blocchi.
- Supporto a `config.ini` per ordinare, nascondere colonne e normalizzare campi data.
- Pagina HTML generata con tabella ordinabile/filtrabile e vista a schede.
- Normalizzazione automatica del campo "Time Played" (secondi → formato `h:mm`, ordinabile numericamente).
- Upload opzionale del file generato su server WebDAV tramite flag `--upload`.
- Script `lista.sh` per l'esecuzione automatizzata (creazione venv, generazione pagina).
- Licenza MIT.

### Changed
- Codice riorganizzato da un singolo script monolitico a un pacchetto Python
  (`playnite_page/`) con moduli separati per parsing, configurazione, rendering e upload WebDAV.
- Corretta incoerenza tra le variabili d'ambiente usate dal codice WebDAV (`WEBDAV_USER`/`WEBDAV_PASS`)
  e quelle documentate (`WEBDAV_USERNAME`/`WEBDAV_PASSWORD`) — ora allineate.
- `lista.sh` aggiornato per creare il virtualenv solo se assente e installare le dipendenze automaticamente.
- README riscritto con struttura del progetto, tabella parametri e istruzioni di configurazione aggiornate.
