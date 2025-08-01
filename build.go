package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
)

// static page info struct declaration
var pages = []struct {
	File 		string
	Title		string
	Tmpl 		string
	Output 	string
} {
		{"index.html", "Home", "index.html", "docs/index.html"},
		{"about.html", "About Me", "about.html", "docs/about.html"},
		{"projects.html", "Projects", "projects.html", "docs/projects.html"},
		{"contact.html", "Contact", "contact.html", "docs/contact.html"},
		{"404.html", "Not Found", "404.html", "docs/404.html"},
}

func main() {
	err := os.MkdirAll("docs", 0755)
	if err != nil {
		log.Fatal("Failed to create docs directory:", err)
	}

	for _, page := range pages {
		f, err := os.Create(page.Output)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v", page.Output, err)
		}
		defer f.Close()

		tmpls, err := template.ParseFiles(
			filepath.Join("templates", "base.html"),
			filepath.Join("templates", page.Tmpl),
		)
		if err != nil {
				log.Fatal("Failed to parse templates:", err)
		}

		err = tmpls.ExecuteTemplate(f, page.Tmpl, map[string]any{
			"Title": page.Title,
		})
		if err != nil {
			log.Fatalf("Failed to render %s: %v", page.File, err)
		}
		log.Printf("Rendered %s → %s", page.Tmpl, page.Output)
	}

	// Copy static files
	copyStatic("static", "docs/static")
}

// Copies static/ → docs/static/
func copyStatic(src, dest string) {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		log.Fatal(err)
	}

	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		destPath := filepath.Join(dest, rel)

		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		in, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, in, 0644)
	})
}
