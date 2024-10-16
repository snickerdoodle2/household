package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
)

var value bool

func main() {
	port := flag.Int("port", 8888, "Server port")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", statusHandler)
	mux.HandleFunc("GET /value", getValueHandler)
	mux.HandleFunc("PUT /value", setValueHandler)
	mux.HandleFunc("POST /value", toggleValueHandler)

	value = false

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server listening on http://localhost:%d\n", *port)
	http.ListenAndServe(addr, mux)
}

func logger(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request", "url", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{ \"status\": \"online\", \"type\": \"binary_switch\" }"))
}

func getValueHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"value\": %t }", value)
}

func setValueHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Value bool `json:"value"`
	}

	maxBytes := 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&input)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{ \"error\": \"Field value should be boolean\" }")
		return
	}

	value = input.Value
	fmt.Printf("Got: %t\n", value)
	fmt.Fprintf(w, "{ \"value\": %t }", value)
}

func toggleValueHandler(w http.ResponseWriter, r *http.Request) {
	value = !value
	fmt.Fprintf(w, "{ \"value\": %t }", value)
}
