package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	webDir := filepath.Join(cwd, "web")
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Web directory not found: %s", webDir)
	}

	fs := http.FileServer(http.Dir(webDir))

	// Reverse proxy to backend for API requests (avoids CORS)
	backendURL, _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		r.Host = backendURL.Host
		r.URL.Scheme = backendURL.Scheme
		r.URL.Host = backendURL.Host
		proxy.ServeHTTP(w, r)
	}

	http.HandleFunc("/api/", apiHandler)

	// Static and SPA handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(webDir, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/js/") || strings.HasPrefix(r.URL.Path, "/css/") {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	// Determine port
	port := os.Getenv("FRONTEND_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "3000"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("Frontend server starting on http://localhost%s", port)
	log.Printf("Serving files from: %s", webDir)
	log.Printf("Proxying API to: %s", backendURL.String())

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
