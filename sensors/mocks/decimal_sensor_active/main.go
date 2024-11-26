package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Config struct {
	minValue float64
	maxValue float64
}

type Measurement struct {
	MessageType string  `json:"message-type"`
	SensorType  string  `json:"sensor-type"`
	Value       float64 `json:"value"`
	IdToken     string  `json:"id-token"`
}

var (
	cfg                        Config
	serverUri                  string
	serverMeasurementsEndpoint string
	serverInitAckEndpoint      string
	idToken                    string
	initWg                     sync.WaitGroup
	interval                   *int
)

func main() {
	port := flag.Int("port", 8888, "Server port")
	flag.Float64Var(&cfg.minValue, "min", 0.0, "Minimal value")
	flag.Float64Var(&cfg.maxValue, "max", 1.0, "Maximal value")
	interval = flag.Int("interval", 5, "Interval between measurements in seconds")
	flag.Parse()

	log.Printf("Starting sensor simulator - waiting for init")
	initWg.Add(1)

	go activeLoop()

	mux := http.NewServeMux()

	mux.HandleFunc("/status", statusHandler)
	mux.HandleFunc("/value", valueHandler)
	mux.HandleFunc("/init", initHandler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func activeLoop() {
	initWg.Wait()

	log.Printf("Initialized - sending measurements to %s every %d seconds\n", serverUri, *interval)

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		utime := time.Now().Unix()
		value := (math.Sin(2*float64(utime))+math.Sin(math.Pi*float64(utime))+2)*(cfg.maxValue-cfg.minValue)/4 + cfg.minValue

		measurement := Measurement{
			MessageType: "measurement",
			SensorType:  "decimal_sensor",
			Value:       value,
			IdToken:     idToken,
		}

		jsonData, err := json.Marshal(measurement)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			continue
		}

		// Ensure proper URL construction
		url := serverUri
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url + serverMeasurementsEndpoint
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error sending request: %v", err)
			continue
		}
		resp.Body.Close()

		log.Printf("Sent measurement: %+v", measurement)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"online","type":"decimal_sensor"}`))
}

func valueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	utime := time.Now().Unix()
	value := (math.Sin(2*float64(utime))+math.Sin(math.Pi*float64(utime))+2)*(cfg.maxValue-cfg.minValue)/4 + cfg.minValue

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"value":%f}`, value)
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Init request received")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		IdToken              string `json:"id-token"`
		ServerUri            string `json:"server-uri"`
		InitAckEndpoint      string `json:"init-ack-endpoint"`
		MeasurementsEndpoint string `json:"measurements-endpoint"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	serverUri = input.ServerUri
	idToken = input.IdToken
	serverMeasurementsEndpoint = input.MeasurementsEndpoint
	serverInitAckEndpoint = input.InitAckEndpoint

	w.WriteHeader(http.StatusOK)

	url := serverUri
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url + serverInitAckEndpoint
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}

	if resp.StatusCode < 300 {
		resp.Body.Close()
		initWg.Done()
	}

}
