package cli

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/example/mycli/internal/config"
	"github.com/example/mycli/internal/report"
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
  mycli run [flags]

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

	targetURL := cfg.LoadTest.Domain + cfg.LoadTest.Endpoint

	fmt.Fprintf(c.stderr, "Starting load test...\n")
	fmt.Fprintf(c.stderr, "Target: %s\n", targetURL)
	fmt.Fprintf(c.stderr, "RPS: %d\n", cfg.LoadTest.RPS)
	fmt.Fprintf(c.stderr, "Concurrency: %d\n", cfg.LoadTest.Concurrency)
	fmt.Fprintf(c.stderr, "Duration: %ds\n", cfg.LoadTest.Duration)
	fmt.Fprintf(c.stderr, "\n")

	// Run load test
	results := c.executeLoadTest(targetURL, cfg.LoadTest.RPS, cfg.LoadTest.Concurrency, cfg.LoadTest.Duration)

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

func (c *CLI) executeLoadTest(url string, rps, concurrency, duration int) *report.Results {
	results := &report.Results{
		URL:         url,
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

	// Channel to distribute work
	workChan := make(chan int, totalRequests)

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range workChan {
				start := time.Now()
				resp, err := client.Get(url)
				elapsed := time.Since(start)

				result := report.RequestResult{
					Timestamp:  start,
					Duration:   elapsed,
					StatusCode: 0,
					Error:      "",
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
				workChan <- requestCount
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
