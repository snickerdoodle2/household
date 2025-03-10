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

type input struct {
	IdToken              string `json:"id-token"`
	ServerUri            string `json:"server-uri"`
	InitAckEndpoint      string `json:"init-ack-endpoint"`
	MeasurementsEndpoint string `json:"measurements-endpoint"`
}

var (
	cfg                        Config
	serverUri                  string
	serverMeasurementsEndpoint string
	serverInitAckEndpoint      string
	idToken                    string
	initWg                     sync.WaitGroup
	interval                   *int
	configured                 bool
	configMutex                sync.Mutex
)

func main() {
	port := flag.Int("port", 8888, "Server port")
	flag.Float64Var(&cfg.minValue, "min", 0.0, "Minimal value")
	flag.Float64Var(&cfg.maxValue, "max", 1.0, "Maximal value")
	interval = flag.Int("interval", 5, "Interval between measurements in seconds")
	flag.Parse()

	log.Printf("Starting sensor - waiting for init")
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

		configMutex.Lock()
		url := serverUri
		measurementsEndpoint := serverMeasurementsEndpoint
		locIdToken := idToken
		configMutex.Unlock()

		measurement := Measurement{
			MessageType: "measurement",
			SensorType:  "decimal_sensor",
			Value:       value,
			IdToken:     locIdToken,
		}

		jsonData, err := json.Marshal(measurement)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			continue
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url + measurementsEndpoint
		}

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		log.Printf("Sending measurement to: %s", url)

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error sending request: %v", err)
			continue
		}

		if resp.StatusCode >= 300 {
			log.Printf("Server returned error status: %d", resp.StatusCode)
		}
		resp.Body.Close()

	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Status request received")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"online","type":"decimal_sensor"}`))
}

func valueHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Value request received")
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := input{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	configMutex.Lock()
	serverUri = input.ServerUri
	idToken = input.IdToken
	serverMeasurementsEndpoint = input.MeasurementsEndpoint
	serverInitAckEndpoint = input.InitAckEndpoint
	configMutex.Unlock()

	log.Printf("Init request received: %+v", input)

	w.WriteHeader(http.StatusOK)

	err := sendInitAck(input)
	if err != nil {
		log.Printf("Error sending init ack: %v", err)
		return
	}
	if !configured {
		initWg.Done()
	}
	configured = true

}

func sendInitAck(input input) error {
	client := &http.Client{}
	client.Timeout = 5 * time.Second

	configMutex.Lock()
	url := serverUri
	initAckEndpoint := serverInitAckEndpoint
	configMutex.Unlock()

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(input)
	if err != nil {
		log.Printf("sendInitAck marshall error: %s", err.Error())
		return err
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url + initAckEndpoint
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Printf("handleRuleRequests request creation error: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	log.Printf("Sending init ack to %s", url)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("handleInitRequest request error: %s", err.Error())
		return err
	}

	if resp.StatusCode < 300 {
		resp.Body.Close()
		return nil
	} else {
		log.Printf("init response code: %d", resp.StatusCode)
		return fmt.Errorf("init response code: %d", resp.StatusCode)
	}
}
