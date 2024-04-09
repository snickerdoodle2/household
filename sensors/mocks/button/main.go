package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

var value bool
var debounceTime time.Time
var delay int

func main() {
	port := flag.Int("port", 8888, "Server port")
	flag.IntVar(&delay, "delay", 2, "Debounce delay in seconds")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", statusHandler)
	mux.HandleFunc("GET /value", getValueHandler)
	mux.HandleFunc("POST /value", pushHandler)
	value = false

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server listening on http://localhost:%d", *port)
	http.ListenAndServe(addr, mux)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{ \"status\": \"online\", \"type\": \"button\" }"))
}

func getValueHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"value\": %t }", value)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	value = true

	go func() {
		debounceTime = time.Now().Add(time.Duration(delay) * time.Second)
		time.Sleep(time.Duration(delay) * time.Second)
		if time.Now().After(debounceTime) {
			value = false
		}
	}()

	fmt.Fprintf(w, "{ \"value\": %t }", value)
}
