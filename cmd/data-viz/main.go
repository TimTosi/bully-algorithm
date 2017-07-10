package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// -----------------------------------------------------------------------------

var ds *DataServer

// -----------------------------------------------------------------------------

// homeHandler returns the data-viz UI.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), 302)
	}
	tmpl.ExecuteTemplate(w, tmpl.Name(), nil)
}

// dataHandler returns a JSON representing current nodes state.
func dataHandler(w http.ResponseWriter, r *http.Request) {
	dataSet, err := json.Marshal(ds.Data())
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataSet)
}

// -----------------------------------------------------------------------------

func main() {
	var err error
	if ds, err = NewDataServer("0.0.0.0:8081"); err != nil {
		log.Fatal(err)
	}
	go ds.Listen()

	http.Handle("/assets/", http.StripPrefix(
		"/assets/",
		http.FileServer(http.Dir("assets")),
	))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", dataHandler)
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
