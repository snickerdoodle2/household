package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"time"
)

type Config struct {
	minValue float64
	maxValue float64
}

var cfg Config

func main() {
	port := flag.Int("port", 8888, "Server port")
	flag.Float64Var(&cfg.minValue, "min", 0.0, "Minimal value")
	flag.Float64Var(&cfg.maxValue, "max", 1.0, "Maximal value")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /status", statusHandler)
	mux.HandleFunc("GET /value", valueHandler)

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server listening on http://localhost:%d", *port)
	http.ListenAndServe(addr, mux)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{ status: \"online\", type: \"decimal_sensor\" }"))
}

func valueHandler(w http.ResponseWriter, r *http.Request) {
	utime := time.Now().Unix()
	value := (math.Sin(2*float64(utime))+math.Sin(math.Pi*float64(utime))+2)*(cfg.maxValue-cfg.minValue)/4 + cfg.minValue

	fmt.Fprintf(w, "{ value: %f }", value)
}
