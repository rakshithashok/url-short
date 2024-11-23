package main

import (
	"fmt"
	"net/http"
)

// Create a DS for mapping between shortened version and normal URL

// Generate a random short code
func generateShortCode() string {

}

// Shorten URL handler
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {

	// Parse original URL from request body

	// Generate unique short code

	// Store the mapping

	// Send back the short URL as response

}

// Redirect handler
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	// read the request, access your DS to find the full version
	// of URL, redirect to that URL
	shortCode := strings.TrimPrefix(r.URL.Path, "/r/")
	originalURL, exists := urlMap[shortCode]
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}


	http.Redirect(w, r, originalURL, http.StatusFound)

}

// Serve frontend index.html file
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,"index.html")
}

func main() {

	// Route for serving the frontend page
	http.HandleFunc("/", indexHandler)

	// Route for the API to shorten URLs
	http.HandleFunc("/shorten", shortenURLHandler)

	// Route for handling redirects
	http.HandleFunc("/r/", redirectHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
