package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
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

var isDevMode bool = true

// Global map to store pre-parsed template sets
var templates = make(map[string]*template.Template)

func main() {
	// check for env variable to set dev mode for deployment
	if os.Getenv("ENV") == "production" {
		isDevMode = false
		log.Println("Running in PRODUCTION mode. Templates will be pre-parsed.")
	} else {
		log.Println("Running in DEVELOPMENT mode. Templates will be re-parsed on each request.")
	}

	// pre-parse templates on server startup
	if !isDevMode {
		if err := parseTemplates(); err != nil {
			log.Fatalf("Error parsing templates: %v", err)
		}
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/projects", projectsHandler)
	mux.HandleFunc("/experiments", experimentsHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	if !isDevMode {
		log.Println("Server starting on :3000")
		log.Fatal(http.ListenAndServe(":3000", mux))
	} else {
		log.Println("Server starting on :8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
	}
}

// only called in production
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

func renderTemplate(w http.ResponseWriter, name string, data any) {
	var tmpl *template.Template
	var err error
	if isDevMode {
		// in Dev mode, re-parse templates on each request
		log.Printf("Dev mode: Re-parsing template '%s'...\n", name)
		basePath := filepath.Join("templates", "base.html")
		tmplPath := filepath.Join("templates", name+".html")

		tmpl, err = template.ParseFiles(basePath, tmplPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Dev mode render error: %v\n", err)
			return
		}
	} else {
		var ok bool
		tmpl, ok = templates[name]
		if !ok {
			http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
			log.Printf("Error: Template '%s' not found in pre-parsed map.\n", name)
			return
		}
	}
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Render error for %s: %v", name, err)
	}
}

//	func rootHandler(w http.ResponseWriter, r *http.Request) {
//		http.Redirect(w, r, "/home", http.StatusFound)
//	}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r) // If not root and no other route matches, serve 404
		return
	}
	p := Page{Title: "Home"}
	renderTemplate(w, "index", &p)
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "About"}
	renderTemplate(w, "about", &p)
}
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "Projects"}
	renderTemplate(w, "projects", &p)
}
func experimentsHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "Experiments"}
	renderTemplate(w, "experiments", &p)
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "404 Not Found"}
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(w, "404", &p)
}
