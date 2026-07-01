package main

import (
	"bufio"
	"os"
	"strings"
)

// Config rispecchia la sezione [columns] di config.ini.
// In Go i campi che iniziano con maiuscola sono "esportati" (pubblici),
// visibili da altri package — qui non ci serve perché siamo tutti in main,
// ma è la convenzione che troverai ovunque.
type Config struct {
	Order []string
	Hide  []string
	Dates []string
}

// loadConfig legge config.ini. Se path è vuoto o il file non esiste,
// ritorna una Config vuota senza errore: in Python facevamo lo stesso
// con `if not path or not path.exists(): return empty`.
func loadConfig(path string) (Config, error) {
	var cfg Config

	if path == "" {
		return cfg, nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close() // defer = "esegui questo alla fine della funzione", qualunque sia il percorso di uscita

	inColumnsSection := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inColumnsSection = strings.EqualFold(line, "[columns]")
			continue
		}
		if !inColumnsSection {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		values := splitCSVList(value)

		switch strings.ToLower(key) {
		case "order":
			cfg.Order = values
		case "hide":
			cfg.Hide = values
		case "dates":
			cfg.Dates = values
		}
	}

	// scanner.Err() ritorna un errore diverso da nil se qualcosa è andato
	// storto DURANTE la scansione (non a fine file, quello è normale).
	return cfg, scanner.Err()
}

// splitCSVList spezza "A,B, C" in ["A","B","C"], scartando spazi ed elementi vuoti.
func splitCSVList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

// contains verifica se slice contiene value.
// Nota: da Go 1.21 esiste slices.Contains nella libreria standard,
// la scriviamo a mano qui solo per esercizio.
func contains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}
