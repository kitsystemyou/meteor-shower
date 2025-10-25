package report

import (
	"encoding/json"
	"fmt"
	"io"
)

type JSONReport struct {
	URL          string                 `json:"url"`
	RPS          int                    `json:"rps"`
	Concurrency  int                    `json:"concurrency"`
	Duration     int                    `json:"duration"`
	StartTime    string                 `json:"start_time"`
	EndTime      string                 `json:"end_time"`
	Statistics   JSONStatistics         `json:"statistics"`
	StatusCodes  map[int]int            `json:"status_codes"`
	Requests     []JSONRequestResult    `json:"requests,omitempty"`
}

type JSONStatistics struct {
	TotalRequests   int     `json:"total_requests"`
	SuccessRequests int     `json:"success_requests"`
	FailedRequests  int     `json:"failed_requests"`
	TotalDurationMs int64   `json:"total_duration_ms"`
	MinDurationMs   int64   `json:"min_duration_ms"`
	MaxDurationMs   int64   `json:"max_duration_ms"`
	AvgDurationMs   int64   `json:"avg_duration_ms"`
	MedianDurationMs int64  `json:"median_duration_ms"`
	P95DurationMs   int64   `json:"p95_duration_ms"`
	P99DurationMs   int64   `json:"p99_duration_ms"`
	RequestsPerSec  float64 `json:"requests_per_sec"`
}

type JSONRequestResult struct {
	Timestamp   string `json:"timestamp"`
	DurationMs  int64  `json:"duration_ms"`
	StatusCode  int    `json:"status_code"`
	Error       string `json:"error,omitempty"`
}

func GenerateJSON(w io.Writer, results *Results) error {
	stats := results.CalculateStatistics()

	report := JSONReport{
		URL:         results.URL,
		RPS:         results.RPS,
		Concurrency: results.Concurrency,
		Duration:    results.Duration,
		StartTime:   results.StartTime.Format("2006-01-02T15:04:05Z07:00"),
		EndTime:     results.EndTime.Format("2006-01-02T15:04:05Z07:00"),
		Statistics: JSONStatistics{
			TotalRequests:    stats.TotalRequests,
			SuccessRequests:  stats.SuccessRequests,
			FailedRequests:   stats.FailedRequests,
			TotalDurationMs:  stats.TotalDuration.Milliseconds(),
			MinDurationMs:    stats.MinDuration.Milliseconds(),
			MaxDurationMs:    stats.MaxDuration.Milliseconds(),
			AvgDurationMs:    stats.AvgDuration.Milliseconds(),
			MedianDurationMs: stats.MedianDuration.Milliseconds(),
			P95DurationMs:    stats.P95Duration.Milliseconds(),
			P99DurationMs:    stats.P99Duration.Milliseconds(),
			RequestsPerSec:   stats.RequestsPerSec,
		},
		StatusCodes: stats.StatusCodeCounts,
		Requests:    make([]JSONRequestResult, 0, len(results.Requests)),
	}

	// Include individual request results
	for _, req := range results.Requests {
		report.Requests = append(report.Requests, JSONRequestResult{
			Timestamp:  req.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			DurationMs: req.Duration.Milliseconds(),
			StatusCode: req.StatusCode,
			Error:      req.Error,
		})
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
