package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var ds *DataServer

// homeHandler returns the data-viz UI.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusFound)
	}
	if err := tmpl.ExecuteTemplate(w, tmpl.Name(), nil); err != nil {
		log.Fatalf("homeHandler: %+v", err)
	}
}

// dataHandler returns a JSON representing current nodes state.
func dataHandler(w http.ResponseWriter, r *http.Request) {
	dataSet, err := json.Marshal(ds.Data())
	if err != nil {
		log.Fatalf("dataHandler: %+v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(dataSet)
}

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
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("data-viz: %+v", err)
	}
}
