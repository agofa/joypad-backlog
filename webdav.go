package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// loadDotEnv legge un file .env molto semplice (KEY=VALUE per riga) e
// imposta le variabili d'ambiente di conseguenza, solo se non già presenti.
// python-dotenv fa lo stesso lavoro in una riga (load_dotenv()); qui lo
// scriviamo a mano perché è un buon esercizio di I/O testuale in Go
// e ci evita una dipendenza esterna per una cosa così piccola.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // .env assente: non è un errore, semplicemente non c'è nulla da caricare
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// uploadWebDAV carica il file generato sul server WebDAV configurato
// tramite variabili d'ambiente (o file .env nella cartella corrente).
func uploadWebDAV(filePath string) error {
	loadDotEnv(".env")

	hostname := os.Getenv("WEBDAV_URL")
	username := os.Getenv("WEBDAV_USERNAME")
	password := os.Getenv("WEBDAV_PASSWORD")

	var missing []string
	for name, val := range map[string]string{
		"WEBDAV_URL":      hostname,
		"WEBDAV_USERNAME": username,
		"WEBDAV_PASSWORD": password,
	} {
		if val == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("variabili d'ambiente mancanti: %s", strings.Join(missing, ", "))
	}

	remotePath := os.Getenv("WEBDAV_PATH")
	if remotePath == "" {
		remotePath = filepath.Base(filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	url := strings.TrimRight(hostname, "/") + "/" + strings.TrimLeft(remotePath, "/")

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "text/html")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// WebDAV risponde 201 (Created) per un file nuovo, 204 (No Content) per un overwrite.
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("il server ha risposto con status %d", resp.StatusCode)
	}

	fmt.Println("✅ Upload completato su WebDAV:", remotePath)
	return nil
}
