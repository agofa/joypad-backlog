package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Row rappresenta un gioco: in Python era un dict, qui è una map.
// map[string]string è l'equivalente diretto di Dict[str, str].
type Row map[string]string

// parseCSV legge il file di input e ritorna (righe, intestazioni).
// In Go le funzioni possono ritornare più valori — qui ne ritorniamo tre,
// col terzo che è sempre `error`: è LA convenzione Go per la gestione errori
// (niente eccezioni/try-except come in Python).
func parseCSV(path string) ([]Row, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("apertura file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true // tollera virgolette non perfettamente escaped, come il csv.DictReader di Python

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("lettura csv: %w", err)
	}
	if len(records) == 0 {
		return nil, nil, nil
	}

	headers := records[0]
	// Il file esportato da Playnite ha un BOM UTF-8 davanti al primo header.
	// Il csv.Reader di Go, a differenza di utf-8-sig in Python, non lo rimuove da solo.
	headers[0] = strings.TrimPrefix(headers[0], "\ufeff")

	rows := make([]Row, 0, len(records)-1)
	for _, record := range records[1:] {
		row := make(Row, len(headers))
		for i, h := range headers {
			if i < len(record) {
				row[h] = record[i]
			}
		}
		rows = append(rows, row)
	}

	return rows, headers, nil
}

// parseDate prova gg/mm/aaaa [hh:mm:ss]. Ritorna (data, ok) invece di
// sollevare un errore: in Go è idiomatico per i casi "non trovato" usare
// un secondo valore booleano, come fa anche time.Parse in modo diverso.
func parseDate(s string) (time.Time, bool) {
	for _, layout := range []string{"02/01/2006 15:04:05", "02/01/2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// formatSeconds converte secondi totali in "h:mm".
func formatSeconds(totalSeconds int) string {
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	return fmt.Sprintf("%d:%02d", hours, minutes)
}

// secondsFromRaw converte la stringa grezza di "Time Played" in int,
// tollerante a valori vuoti o non numerici (torna 0 in quel caso).
func secondsFromRaw(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	// Playnite esporta a volte un float ("123.0"): passiamo da float64 e poi troncamo,
	// esattamente come faceva int(float(raw)) in Python.
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return int(f)
}
