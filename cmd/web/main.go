package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Set the web directory path
	webDir := filepath.Join(cwd, "web")

	// Check if web directory exists
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		log.Fatalf("Web directory not found: %s", webDir)
	}

	// Create a file server for the web directory
	fs := http.FileServer(http.Dir(webDir))

	// Handle all routes by serving the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the requested file
		filePath := filepath.Join(webDir, r.URL.Path)
		
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// If file doesn't exist, serve index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
			return
		}
		
		// Serve the file
		fs.ServeHTTP(w, r)
	})

	// Add CORS headers for development
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Proxy API requests to the main server
		http.Redirect(w, r, "http://localhost:8080"+r.URL.Path, http.StatusTemporaryRedirect)
	})

	port := ":3000"
	log.Printf("Frontend server starting on http://localhost%s", port)
	log.Printf("Serving files from: %s", webDir)
	log.Printf("API requests will be proxied to: http://localhost:8080")
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
