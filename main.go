package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, name string, data any)  {
	log.Println("Template rendering...")
	template, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", name),
		)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = template.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Render error: %v", err)
	}
}

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index.html", map[string]string{"Title": "Home"})
	})
	http.HandleFunc("/about", func (w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "about.html", map[string]string{"Title": "About"})
	})
	http.HandleFunc("/contact", func (w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "contact.html", map[string]string{"Title": "Contact Info"})
	})
	http.HandleFunc("/projects", func (w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "projects.html", map[string]string{"Title": "Projects"})
	})
	http.HandleFunc("/404", func (w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "404.html", map[string]string{"Title": "404 Not Found"})
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}


