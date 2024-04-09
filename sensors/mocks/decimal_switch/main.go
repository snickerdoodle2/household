package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

var value float32

func main() {
	port := flag.Int("port", 8888, "Server port")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", statusHandler)
	mux.HandleFunc("GET /value", getValueHandler)
	mux.HandleFunc("PUT /value", setValueHandler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server listening on http://localhost:%d", *port)
	http.ListenAndServe(addr, mux)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{ \"status\": \"online\", \"type\": \"decimal_switch\" }"))
}

func getValueHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"value\": %f }", value)
}

func setValueHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Value float32 `json:"value"`
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
	fmt.Fprintf(w, "{ \"value\": %f }", value)
}
