package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles(
		"./templates/index.html",
		"./templates/base.html",
		// "./templates/about.html",
		// "./templates/404.html",
		// "./templates/contact.html",
		// "./templates/projects.html",
	))

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}


