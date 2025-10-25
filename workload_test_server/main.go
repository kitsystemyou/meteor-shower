package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	port         = flag.Int("port", 8080, "Port to listen on")
	delay        = flag.Int("delay", 10, "Response delay in milliseconds")
	errorRate    = flag.Float64("error-rate", 0.0, "Error rate (0.0 to 1.0)")
	randomDelay  = flag.Bool("random-delay", false, "Add random delay variation")
	requestCount int64
)

type Stats struct {
	TotalRequests int64     `json:"total_requests"`
	Uptime        string    `json:"uptime"`
	StartTime     time.Time `json:"start_time"`
}

var startTime time.Time

func main() {
	flag.Parse()

	startTime = time.Now()
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/stats", handleStats)
	http.HandleFunc("/slow", handleSlow)
	http.HandleFunc("/error", handleError)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting test server on %s", addr)
	log.Printf("Configuration:")
	log.Printf("  - Base delay: %dms", *delay)
	log.Printf("  - Error rate: %.2f%%", *errorRate*100)
	log.Printf("  - Random delay: %v", *randomDelay)
	log.Printf("\nEndpoints:")
	log.Printf("  GET /         - Normal endpoint with configured delay")
	log.Printf("  GET /health   - Health check (no delay)")
	log.Printf("  GET /stats    - Server statistics")
	log.Printf("  GET /slow     - Slow endpoint (500ms delay)")
	log.Printf("  GET /error    - Always returns 500 error")
	log.Printf("\n")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	log.Printf("requested")

	// Simulate error
	if *errorRate > 0 && rand.Float64() < *errorRate {
		http.Error(w, "Random error", http.StatusInternalServerError)
		return
	}

	// Simulate delay
	delayMs := *delay
	if *randomDelay {
		// Add ±50% variation
		variation := int(float64(*delay) * 0.5)
		delayMs = *delay + rand.Intn(variation*2) - variation
		if delayMs < 0 {
			delayMs = 0
		}
	}
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"delay_ms":  delayMs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)
	stats := Stats{
		TotalRequests: atomic.LoadInt64(&requestCount),
		Uptime:        uptime.String(),
		StartTime:     startTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleSlow(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	time.Sleep(500 * time.Millisecond)

	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"delay_ms":  500,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleError(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	http.Error(w, "Intentional error", http.StatusInternalServerError)
}
