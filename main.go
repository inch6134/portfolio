package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Page struct { // Page parameters to pass to template
	Title string
	// ProjectData *Project
	// IsDarkMode bool

	// Add other dynamic fields here, e.g.,
	// Content string
	// Data map[string]interface{}
	// Items []string
}

// type Project struct {
// 	Name string
// 	RepoURL string
// 	Description template.HTML
// 	Technologies []string
// } 

// Global map to store pre-parsed template sets
var templates = make(map[string]*template.Template)

func main() {
	// pre-parse templates on server startup
	if err := parseTemplates(); err != nil {
		log.Fatalf("Error parsing templates", err)
	}
	
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/404", fourohfourHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseTemplates() error {
	baseTemplatePath := filepath.Join("templates", "base.html")

	contentTemplates := []string{
		"index.html",
		"about.html",
		"projects.html",
		"404.html",
	}

	for _, tmplFile := range contentTemplates {
		tmplPath := filepath.Join("templates", tmplFile)
		// template name w/o extension is key for the map
		tmplName := tmplFile[:len(tmplFile)-len(".html")]

		// create set of new templates for each page
		tmpl, err := template.ParseFiles(baseTemplatePath, tmplPath)
		if err != nil {
			return err
		}
		templates[tmplName] = tmpl
	}
	return nil
}

func renderTemplate(w http.ResponseWriter, name string, data any)  {
	tmpl, ok := templates[name]
	if !ok{
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		log.Printf("Error: Template '%s' not found in pre-parsed map.\n", name)
		return
	}
	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Render error for %s: %v", name, err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"Home"}
	renderTemplate(w, "index", &p)
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"About"}
	renderTemplate(w, "about", &p)
}
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"Projects"}
	renderTemplate(w, "projects", &p)
}
func fourohfourHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title:"404 Not Found"}
	renderTemplate(w, "404", &p)
}




