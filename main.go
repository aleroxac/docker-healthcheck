package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "UP"}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("GET /status", statusHandler)

	api_port := os.Getenv("API_PORT")
	fmt.Printf("Listening on port %s\n", api_port)
	http.ListenAndServe(fmt.Sprintf(":%s", api_port), nil)
}
