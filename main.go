package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatal("Error parsing templates:", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler("index.html", "Home"))
	mux.HandleFunc("/about", makeHandler("about.html", "About Me"))
	mux.HandleFunc("/projects", makeHandler("projects.html", "Projects"))
	mux.HandleFunc("/contact", makeHandler("contact.html", "Contact"))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/resume", http.FileServer(http.Dir("static"))) // /resume will serve static/resume.pdf

	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/404", notFoundHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", with404Handler(mux)))
}

func makeHandler(templateName, title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, templateName, map[string]interface{}{"Title": title})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.ExecuteTemplate(w, "404.html", map[string]interface{}{
		"Title": "Page Not Found",
	})
}

// Middleware: custom 404 handler for undefined routes
func with404Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)
		if rec.statusCode == http.StatusNotFound {
			notFoundHandler(w, r)
		}
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

