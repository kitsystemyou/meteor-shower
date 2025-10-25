package report

import (
	"sort"
	"time"
)

type Results struct {
	URLs        []string
	RPS         int
	Concurrency int
	Duration    int
	StartTime   time.Time
	EndTime     time.Time
	Requests    []RequestResult
}

type RequestResult struct {
	Timestamp  time.Time
	Duration   time.Duration
	StatusCode int
	Error      string
	URL        string
}

type Statistics struct {
	TotalRequests    int
	SuccessRequests  int
	FailedRequests   int
	TotalDuration    time.Duration
	MinDuration      time.Duration
	MaxDuration      time.Duration
	AvgDuration      time.Duration
	MedianDuration   time.Duration
	P95Duration      time.Duration
	P99Duration      time.Duration
	RequestsPerSec   float64
	StatusCodeCounts map[int]int
	URLCounts        map[string]int
}

func (r *Results) CalculateStatistics() Statistics {
	stats := Statistics{
		TotalRequests:    len(r.Requests),
		StatusCodeCounts: make(map[int]int),
		URLCounts:        make(map[string]int),
	}

	if stats.TotalRequests == 0 {
		return stats
	}

	durations := make([]time.Duration, 0, stats.TotalRequests)
	var totalDuration time.Duration

	for _, req := range r.Requests {
		if req.Error == "" {
			stats.SuccessRequests++
			stats.StatusCodeCounts[req.StatusCode]++
		} else {
			stats.FailedRequests++
		}

		if req.URL != "" {
			stats.URLCounts[req.URL]++
		}

		durations = append(durations, req.Duration)
		totalDuration += req.Duration
	}

	// Sort durations for percentile calculations
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	stats.TotalDuration = r.EndTime.Sub(r.StartTime)
	stats.MinDuration = durations[0]
	stats.MaxDuration = durations[len(durations)-1]
	stats.AvgDuration = totalDuration / time.Duration(stats.TotalRequests)
	stats.MedianDuration = durations[len(durations)/2]
	stats.P95Duration = durations[int(float64(len(durations))*0.95)]
	stats.P99Duration = durations[int(float64(len(durations))*0.99)]
	stats.RequestsPerSec = float64(stats.TotalRequests) / stats.TotalDuration.Seconds()

	return stats
}
