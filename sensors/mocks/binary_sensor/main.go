package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"time"
)

func main() {
	port := flag.Int("port", 8888, "Server port")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", statusHandler)
	mux.HandleFunc("GET /value", valueHandler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server listening on http://localhost:%d", *port)
	http.ListenAndServe(addr, mux)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{ \"status\": \"online\", \"type\": \"binary_sensor\" }"))
}

func valueHandler(w http.ResponseWriter, r *http.Request) {
	utime := time.Now().Unix()
	value := math.Sin(2*float64(utime)) + math.Sin(math.Pi*float64(utime))
	var res float64
	if value > 0 {
		res = 1
	} else {
		res = 0
	}

	fmt.Fprintf(w, "{ \"value\": %f }", res)
}
