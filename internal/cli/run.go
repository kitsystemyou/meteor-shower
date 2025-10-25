package cli

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/kitsystemyou/meteor-shower/internal/config"
	"github.com/kitsystemyou/meteor-shower/internal/report"
)

func (c *CLI) runCommand(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	configFile := fs.String("config", "", "config file (default is ./config.yaml)")
	rps := fs.Int("rps", 0, "requests per second (overrides config)")
	concurrency := fs.Int("concurrency", 0, "number of concurrent clients (overrides config)")
	output := fs.String("output", "", "output format: html, json (overrides config)")
	outputShort := fs.String("o", "", "output format: html, json (overrides config)")

	fs.Usage = func() {
		usage := `Run executes load test against the target endpoint.

Usage:
  meteor-shower run [flags]

Flags:
  --rps int              requests per second (overrides config)
  --concurrency int      number of concurrent clients (overrides config)
  -o, --output string    output format: html, json (overrides config)

Global Flags:
  --config string   config file (default is ./config.yaml)
`
		fmt.Fprint(c.stderr, usage)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with command-line flags
	if *rps > 0 {
		cfg.LoadTest.RPS = *rps
	}
	if *concurrency > 0 {
		cfg.LoadTest.Concurrency = *concurrency
	}
	if *output != "" {
		cfg.LoadTest.Output = *output
	} else if *outputShort != "" {
		cfg.LoadTest.Output = *outputShort
	}

	// Validate configuration
	if cfg.LoadTest.RPS <= 0 {
		return fmt.Errorf("rps must be greater than 0")
	}
	if cfg.LoadTest.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}
	if cfg.LoadTest.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	if len(cfg.LoadTest.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint must be specified")
	}

	// Normalize endpoint weights (set default to 1.0 if not specified)
	for i := range cfg.LoadTest.Endpoints {
		if cfg.LoadTest.Endpoints[i].Weight <= 0 {
			cfg.LoadTest.Endpoints[i].Weight = 1.0
		}
	}

	// Build target URLs
	var targets []string
	for _, ep := range cfg.LoadTest.Endpoints {
		targets = append(targets, cfg.LoadTest.Domain+ep.Path)
	}

	fmt.Fprintf(c.stderr, "Starting load test...\n")
	fmt.Fprintf(c.stderr, "Domain: %s\n", cfg.LoadTest.Domain)
	fmt.Fprintf(c.stderr, "Endpoints: %d\n", len(targets))
	for i, ep := range cfg.LoadTest.Endpoints {
		fmt.Fprintf(c.stderr, "  [%d] %s (weight: %.2f)\n", i+1, ep.Path, ep.Weight)
	}
	fmt.Fprintf(c.stderr, "RPS: %d\n", cfg.LoadTest.RPS)
	fmt.Fprintf(c.stderr, "Concurrency: %d\n", cfg.LoadTest.Concurrency)
	fmt.Fprintf(c.stderr, "Duration: %ds\n", cfg.LoadTest.Duration)
	fmt.Fprintf(c.stderr, "\n")

	// Run load test
	results := c.executeLoadTest(targets, cfg.LoadTest.Endpoints, cfg.LoadTest.RPS, cfg.LoadTest.Concurrency, cfg.LoadTest.Duration)

	// Generate report
	switch cfg.LoadTest.Output {
	case "json":
		return report.GenerateJSON(c.stdout, results)
	case "html":
		return report.GenerateHTML(c.stdout, results)
	default:
		return fmt.Errorf("unsupported output format: %s", cfg.LoadTest.Output)
	}
}

func (c *CLI) executeLoadTest(urls []string, endpoints []config.Endpoint, rps, concurrency, duration int) *report.Results {
	results := &report.Results{
		URLs:        urls,
		RPS:         rps,
		Concurrency: concurrency,
		Duration:    duration,
		StartTime:   time.Now(),
		Requests:    make([]report.RequestResult, 0),
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Calculate interval between requests
	interval := time.Second / time.Duration(rps)
	totalRequests := rps * duration

	// Create HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Build weighted URL selector
	type weightedURL struct {
		url    string
		weight float64
	}
	weightedURLs := make([]weightedURL, 0)
	totalWeight := 0.0
	for i, ep := range endpoints {
		weight := ep.Weight
		if weight <= 0 {
			weight = 1.0
		}
		weightedURLs = append(weightedURLs, weightedURL{url: urls[i], weight: weight})
		totalWeight += weight
	}

	// Normalize weights
	for i := range weightedURLs {
		weightedURLs[i].weight /= totalWeight
	}

	// Function to select URL based on weight
	selectURL := func() string {
		r := rand.Float64()
		cumulative := 0.0
		for _, wu := range weightedURLs {
			cumulative += wu.weight
			if r <= cumulative {
				return wu.url
			}
		}
		return weightedURLs[len(weightedURLs)-1].url
	}

	// Channel to distribute work
	workChan := make(chan string, totalRequests)

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range workChan {
				start := time.Now()
				resp, err := client.Get(url)
				elapsed := time.Since(start)

				result := report.RequestResult{
					Timestamp:  start,
					Duration:   elapsed,
					StatusCode: 0,
					Error:      "",
					URL:        url,
				}

				if err != nil {
					result.Error = err.Error()
				} else {
					result.StatusCode = resp.StatusCode
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}

				mu.Lock()
				results.Requests = append(results.Requests, result)
				mu.Unlock()
			}
		}()
	}

	// Send requests at specified rate
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	requestCount := 0
	timeout := time.After(time.Duration(duration) * time.Second)

	for {
		select {
		case <-timeout:
			close(workChan)
			wg.Wait()
			results.EndTime = time.Now()
			return results
		case <-ticker.C:
			if requestCount < totalRequests {
				workChan <- selectURL()
				requestCount++
			} else {
				close(workChan)
				wg.Wait()
				results.EndTime = time.Now()
				return results
			}
		}
	}
}


