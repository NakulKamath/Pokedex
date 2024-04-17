package main

import (
	"net/http"
)

func main() {
	// Serve the frontend1 directory
	fs := http.FileServer(http.Dir("C:\\Users\\prana\\OneDrive\\Documents\\college\\pokedex\\frontend1"))
	http.Handle("/", fs)

	// Start the HTTP server
	http.ListenAndServe(":8000", nil)
}
