package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Il package flag è l'equivalente stdlib di argparse.
	// A differenza di Python, qui i flag booleani e le variabili sono
	// tipizzate a compile time: flag.String ritorna *string, non stringa,
	// perché il valore viene scritto DENTRO quel puntatore solo dopo Parse().
	inputPath := flag.String("input", "", "File di input (CSV) — obbligatorio")
	output := flag.String("o", "output.html", "File HTML di output")
	title := flag.String("title", "Pagina generata", "Titolo della pagina HTML")
	templatePath := flag.String("template", "", "Template HTML personalizzato (opzionale)")
	configPath := flag.String("config", "", "File config.ini (opzionale)")
	upload := flag.Bool("upload", false, "Carica il file generato su WebDAV")

	flag.Parse()

	// A differenza di argparse, flag.String non gestisce comodamente un
	// argomento posizionale obbligatorio: lo leggiamo da flag.Args()
	// se non è stato passato con -input, per restare compatibili con
	// `generate_page.py input.csv -o ...` della versione Python.
	if *inputPath == "" {
		if args := flag.Args(); len(args) > 0 {
			*inputPath = args[0]
		}
	}
	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "Errore: specifica il file di input (es. playnite.csv)")
		os.Exit(2)
	}

	if err := run(*inputPath, *output, *title, *templatePath, *configPath, *upload); err != nil {
		fmt.Fprintln(os.Stderr, "Errore:", err)
		os.Exit(1)
	}
}

// run contiene la logica vera e propria, separata da main() per due motivi:
// 1) main() con os.Exit() è scomodo da testare;
// 2) run() ritorna un error "normale", gestito in un unico punto sopra.
func run(inputPath, output, title, templatePath, configPath string, upload bool) error {
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("file non trovato: %s", inputPath)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("lettura config: %w", err)
	}

	rows, headers, err := parseCSV(inputPath)
	if err != nil {
		return fmt.Errorf("lettura csv: %w", err)
	}

	totalSeconds := 0
	for _, row := range rows {
		// Normalizza le colonne data configurate in [columns] dates.
		for _, col := range cfg.Dates {
			if val, ok := row[col]; ok && val != "" {
				if t, ok := parseDate(val); ok {
					row[col] = t.Format("02/01/2006")
				}
			}
		}

		seconds := secondsFromRaw(row["Time Played"])
		totalSeconds += seconds
		row["Time Played__seconds"] = fmt.Sprintf("%d", seconds)
		row["Time Played"] = formatSeconds(seconds)
	}

	sortRowsByName(rows)

	data := PageData{
		Title:       title,
		InputName:   filepath.Base(inputPath),
		SourceType:  "CSV",
		Rows:        rows,
		Columns:     buildColumns(headers, cfg),
		TotalPlayed: formatSeconds(totalSeconds),
		TotalDays:   fmt.Sprintf("%d", totalSeconds/86400),
		Today:       time.Now().Format("02/01/2006 15:04"),
	}

	html, err := renderHTML(data, templatePath)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	if err := os.WriteFile(output, []byte(html), 0644); err != nil {
		return fmt.Errorf("scrittura output: %w", err)
	}
	absPath, _ := filepath.Abs(output)
	fmt.Println("✅ Generato:", absPath)

	if upload {
		if err := uploadWebDAV(output); err != nil {
			return fmt.Errorf("upload webdav: %w", err)
		}
	}

	return nil
}
