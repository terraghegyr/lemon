package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const appsDir = "./apps"

type App struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	// Root files (index.html, style.css)
	http.HandleFunc("/", rootHandler)

	// API
	http.HandleFunc("/api/apps", appsHandler)

	// Apps (secure)
	http.Handle("/apps/",
		http.StripPrefix("/apps/",
			safeFileServer(appsDir),
		),
	)

	log.Println("🍋 lemon running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ---------- Root Files ----------

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.ServeFile(w, r, "static/index.html")
	case "/style.css":
		http.ServeFile(w, r, "static/style.css")
	default:
		http.NotFound(w, r)
	}
}

// ---------- API ----------

func appsHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(appsDir)
	if err != nil {
		http.Error(w, "Failed to read apps", 500)
		return
	}

	var apps []App

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(appsDir, e.Name(), "lemon.json"))
		if err != nil {
			continue
		}

		var app App
		if err := json.Unmarshal(data, &app); err != nil {
			continue
		}

		app.Name = e.Name()
		apps = append(apps, app)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}


func safeFileServer(root string) http.Handler {
	rootAbs, _ := filepath.Abs(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := filepath.Clean(r.URL.Path)
		full := filepath.Join(rootAbs, clean)

		// Prevent directory traversal
		if !strings.HasPrefix(full, rootAbs+string(os.PathSeparator)) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		info, err := os.Stat(full)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Disable directory listing
		if info.IsDir() {
			index := filepath.Join(full, "index.html")
			if _, err := os.Stat(index); err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			full = index
		}

		http.ServeFile(w, r, full)
	})
}

