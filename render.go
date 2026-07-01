package main

import (
	_ "embed" // side-effect import: attiva la direttiva //go:embed qui sotto
	"html/template"
	"os"
	"sort"
	"strings"
)

// go:embed incorpora il contenuto del file nel binario compilato,
// a compile time. Non serve distribuire templates/default.html insieme
// all'eseguibile: è già dentro. In Python questo lo facevamo con una
// stringa Python enorme (DEFAULT_TEMPLATE) proprio per lo stesso motivo.
//
//go:embed templates/default.html
var defaultTemplateSource string

// PageData è tutto ciò che serve al template per renderizzare la pagina.
// In Jinja2 passavamo un dizionario di kwargs; qui, essendo Go tipizzato,
// serve una struct esplicita con un campo per ogni variabile del template.
type PageData struct {
	Title       string
	InputName   string
	SourceType  string
	Rows        []Row
	Columns     []string
	TotalPlayed string
	TotalDays   string
	Today       string
}

// funcMap aggiunge funzioni custom richiamabili dal template.
// html/template non ha un filtro "+1" pronto come Jinja2 (loop.index),
// quindi lo definiamo noi: {{inc $i}} nel template chiama questa funzione.
var funcMap = template.FuncMap{
	"inc": func(i int) int { return i + 1 },
}

// renderHTML compila il template (custom o quello di default) e lo esegue
// con i dati passati, ritornando l'HTML come stringa.
func renderHTML(data PageData, templatePath string) (string, error) {
	var source string
	if templatePath != "" {
		b, err := os.ReadFile(templatePath)
		if err != nil {
			return "", err
		}
		source = string(b)
	} else {
		source = defaultTemplateSource
	}

	// template.New(...).Funcs(...).Parse(...) è un pattern "a catena" molto
	// comune in Go: ogni metodo ritorna l'oggetto stesso (o un puntatore),
	// permettendo di concatenare chiamate senza variabili intermedie.
	tmpl, err := template.New("page").Funcs(funcMap).Parse(source)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// buildColumns applica order/hide da Config per decidere quali colonne
// mostrare e in che ordine. "Name" è sempre esclusa: viene mostrata a parte
// come prima colonna della tabella.
func buildColumns(headers []string, cfg Config) []string {
	hidden := make(map[string]bool, len(cfg.Hide)+1)
	for _, h := range cfg.Hide {
		hidden[h] = true
	}
	hidden["Name"] = true

	var columns []string
	seen := make(map[string]bool)

	for _, col := range cfg.Order {
		if hidden[col] || seen[col] || !contains(headers, col) {
			continue
		}
		columns = append(columns, col)
		seen[col] = true
	}
	for _, h := range headers {
		if hidden[h] || seen[h] {
			continue
		}
		columns = append(columns, h)
		seen[h] = true
	}
	return columns
}

// sortRowsByName ordina le righe per nome, case-insensitive.
// sort.Slice prende una funzione "less" — diverso da Python dove
// rows.sort(key=...) prende una funzione "chiave di ordinamento".
func sortRowsByName(rows []Row) {
	sort.Slice(rows, func(i, j int) bool {
		return strings.ToLower(rows[i]["Name"]) < strings.ToLower(rows[j]["Name"])
	})
}
