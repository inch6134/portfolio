package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Page struct { // Page parameters to pass to template
	Title string
}

func renderTemplate(w http.ResponseWriter, name string, p *Page)  {
	log.Println("Template rendering...")
	template, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", name),
		)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = template.ExecuteTemplate(w, name, p)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Render error: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"Home"}
	renderTemplate(w, "index.html", &p)
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"About"}
	renderTemplate(w, "about.html", &p)
}
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"Projects"}
	renderTemplate(w, "projects.html", &p)
}
func fourohfourHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"404 Not Found"}
	renderTemplate(w, "404.html", &p)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/404", fourohfourHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


