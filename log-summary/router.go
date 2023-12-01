package main

import (
	"html/template"
	"net/http"
	"strconv"
)

func NewRouter() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/add", AddHandler)
	r.HandleFunc("/remove", RemoveHandler)
	r.HandleFunc("/summary", SummaryHandler)
	r.HandleFunc("/image", ImageHandler)

	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	context := struct{ Entries Entries }{AllEntries()}
	template.Must(template.ParseFiles("views/index.html")).Execute(w, context)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	AddEntry(r.PostFormValue("entry"))

	context := struct{ Entries Entries }{AllEntries()}
	template.Must(template.ParseFiles("views/entries.html")).Execute(w, context)
}

func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	RemoveEntry(id)

	context := struct{ Entries Entries }{AllEntries()}
	template.Must(template.ParseFiles("views/entries.html")).Execute(w, context)
}

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	entries := make([]string, 0)
	for _, entry := range AllEntries() {
		entries = append(entries, entry.Content)
	}

	summary := Summarize(entries)

	context := struct{ Summary string }{summary}
	template.Must(template.ParseFiles("views/summary.html")).Execute(w, context)
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	summary := r.PostFormValue("summary")
	image := SummaryImage(summary)

	w.Write([]byte(`
		<body>
		    <turbo-frame id="image">
		        <img src="` + image + `" style="max-width: 100%;">
		    </turbo-frame>
		</body>
	`))
}
