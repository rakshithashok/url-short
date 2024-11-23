package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)
var urlMap = make(map[string]string)
// Create a DS for mapping between shortened version and normal URL

// Generate a random short code
func generateShortCode() string {
    const length = 6
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, length)
    
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    
    return string(b)
}


// Shorten URL handler
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse original URL from request body
	var reqBody struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	originalURL := reqBody.URL
	if originalURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Generate unique short code
	shortCode := generateShortCode()

	// Store the mapping
	urlMap.Lock()
	urlMap.data[shortCode] = originalURL
	urlMap.Unlock()

	// Send back the short URL as response
	shortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

// Redirect handler
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Read the short code from the request URL
	shortCode := strings.TrimPrefix(r.URL.Path, "/r/")
	if shortCode == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	// Access the mapping to find the original URL
	originalURL, exists := urlMap[shortCode]
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
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
