package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type URLMapping struct {
	DefaultURL string
	Conditions map[string]string
}

// For demo purposes, we use an in-memory store protected by RWMutex.
// In production, you'd likely use a persistent database.
var urlStore = struct {
	sync.RWMutex
	m map[string]URLMapping
}{m: make(map[string]URLMapping)}

// generateCode creates a random string of length n which will serve as the short code
func generateCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// Seed the random number generator with the current time
	// rand.Seed(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// createHandler handles POST requests to /create.
// It reads a JSON payload to create a new URLMapping and returns the short URL.
func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON payload into our URLMapping struct.
	var mapping URLMapping
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Generate a unique short code for this URL
	code := generateCode(6)

	// Store the mapping in our in-memory store.
	urlStore.Lock()
	urlStore.m[code] = mapping
	urlStore.Unlock()

	// Build the short URL (assuming the service is hosted on localhost:8080).
	shortUrl := fmt.Sprintf("http://localhost:8080/%s", code)
	resp := map[string]string{"short_url": shortUrl}

	// Return the short URL as JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// redirectHandler handles GET requests for short URLs.
// It determines which destination URL to use based on the request context.
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the short code from the URL path.
	code := strings.TrimPrefix(r.URL.Path, "/")

	// Look up the mapping in our in-memory store.
	urlStore.RLock()
	mapping, exists := urlStore.m[code]
	urlStore.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	// Example condition check: Detect mobile devices by inspecting the User-Agent header.
	ua := strings.ToLower(r.UserAgent())
	var redirectURL string
	if strings.Contains(ua, "mobile") {
		if url, ok := mapping.Conditions["mobile"]; ok {
			redirectURL = url
		}
	}

	// Additional conditions such as geolocation can be added here.

	// If no condition is matched, use the default URL.
	if redirectURL == "" {
		redirectURL = mapping.DefaultURL
	}

	// Redirect the user to the chosen URL.
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("SmartURL server is running on :8080")
	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
