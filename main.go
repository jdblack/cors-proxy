package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {

	parsedArgs := parseArgs("GOCORS", []argument{
		{ Name: "port", Option: "p", Description: "Port to listen on", Default: 8080 },
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers


		// Get the URL from the query parameter
		fullURL := r.URL.Query().Get("url") 
		if fullURL == "" {
			http.Error(w, "URL parameter is required", http.StatusBadRequest)
			return
		}

		fmt.Printf("Received request for URL: %s\n", fullURL)

		// Perform the HTTP GET request to the provided full URL
		resp, err := http.Get(fullURL)
		if err != nil {
			http.Error(w, "Failed to reach the target URL", http.StatusBadGateway)
			log.Printf("Error making request to %s: %v", fullURL, err)
			return
		}
		defer resp.Body.Close()

		// Read the response from the target URL
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			log.Printf("Error reading response: %v", err)
			return
		}


		for k,v := range resp.Header {
			w.Header().Set(k, strings.Join(v," "))
			fmt.Printf("Setting %s %s\n", k, v)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow localhost
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allow specific methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Allow specific headers

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}






		// Write the response from the target URL back to the original response
		log.Printf("Returning status: %v  %v", resp.StatusCode, resp.Status)
		w.WriteHeader(resp.StatusCode)
		w.Write(body)

	})

	// Start the server on port 8080
	log.Printf("LOG: Server starting up on port %v...", parsedArgs["port"])
	fmt.Printf("Server is starting on port %v...", parsedArgs["port"])
	err := http.ListenAndServe(fmt.Sprintf(":%v",parsedArgs["port"]), nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
