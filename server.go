package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

type Server struct {
	client *anthropic.Client
}

func NewServer(client *anthropic.Client) *Server {
	return &Server{client: client}
}

func (s *Server) classifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	 log.Printf("decoded message: content=%q content_type=%q", msg.Content, msg.ContentType)


	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	result, err := classifyWithContext(ctx, s.client, msg)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, result, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, map[string]string{"error": message}, status)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func startServer(client *anthropic.Client) error {
	server := NewServer(client)

	mux := http.NewServeMux()
	mux.HandleFunc("/classify", server.classifyHandler)

	handler := loggingMiddleware(mux)

	log.Printf("Server running on http://localhost:8080")
	return http.ListenAndServe(":8080", handler)
}