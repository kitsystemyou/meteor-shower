package report

import (
	"fmt"
	"html/template"
	"io"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Load Test Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            border-radius: 10px;
            margin-bottom: 20px;
        }
        .header h1 {
            margin: 0 0 10px 0;
        }
        .header p {
            margin: 5px 0;
            opacity: 0.9;
        }
        .section {
            background: white;
            padding: 20px;
            margin-bottom: 20px;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .section h2 {
            margin-top: 0;
            color: #333;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        .stat-card {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        .stat-label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .stat-value {
            font-size: 24px;
            font-weight: bold;
            color: #333;
            margin-top: 5px;
        }
        .status-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 15px;
        }
        .status-table th,
        .status-table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #e0e0e0;
        }
        .status-table th {
            background-color: #f8f9fa;
            font-weight: 600;
            color: #333;
        }
        .status-table tr:hover {
            background-color: #f8f9fa;
        }
        .success {
            color: #28a745;
        }
        .error {
            color: #dc3545;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Load Test Report</h1>
        <p><strong>Target URL:</strong> {{.URL}}</p>
        <p><strong>Test Duration:</strong> {{.Stats.TotalDuration}}</p>
        <p><strong>Start Time:</strong> {{.StartTime.Format "2006-01-02 15:04:05"}}</p>
        <p><strong>End Time:</strong> {{.EndTime.Format "2006-01-02 15:04:05"}}</p>
    </div>

    <div class="section">
        <h2>Configuration</h2>
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-label">Target RPS</div>
                <div class="stat-value">{{.RPS}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Concurrency</div>
                <div class="stat-value">{{.Concurrency}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Duration</div>
                <div class="stat-value">{{.Duration}}s</div>
            </div>
        </div>
    </div>

    <div class="section">
        <h2>Summary</h2>
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-label">Total Requests</div>
                <div class="stat-value">{{.Stats.TotalRequests}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Success</div>
                <div class="stat-value success">{{.Stats.SuccessRequests}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Failed</div>
                <div class="stat-value error">{{.Stats.FailedRequests}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Actual RPS</div>
                <div class="stat-value">{{printf "%.2f" .Stats.RequestsPerSec}}</div>
            </div>
        </div>
    </div>

    <div class="section">
        <h2>Response Time Statistics</h2>
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-label">Min</div>
                <div class="stat-value">{{.Stats.MinDuration}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Average</div>
                <div class="stat-value">{{.Stats.AvgDuration}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Median</div>
                <div class="stat-value">{{.Stats.MedianDuration}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">95th Percentile</div>
                <div class="stat-value">{{.Stats.P95Duration}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">99th Percentile</div>
                <div class="stat-value">{{.Stats.P99Duration}}</div>
            </div>
            <div class="stat-card">
                <div class="stat-label">Max</div>
                <div class="stat-value">{{.Stats.MaxDuration}}</div>
            </div>
        </div>
    </div>

    <div class="section">
        <h2>Status Code Distribution</h2>
        <table class="status-table">
            <thead>
                <tr>
                    <th>Status Code</th>
                    <th>Count</th>
                    <th>Percentage</th>
                </tr>
            </thead>
            <tbody>
                {{range $code, $count := .Stats.StatusCodeCounts}}
                <tr>
                    <td>{{$code}}</td>
                    <td>{{$count}}</td>
                    <td>{{printf "%.2f" (percentage $count $.Stats.TotalRequests)}}%</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>
</html>`

func GenerateHTML(w io.Writer, results *Results) error {
	stats := results.CalculateStatistics()

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"percentage": func(count, total int) float64 {
			if total == 0 {
				return 0
			}
			return float64(count) / float64(total) * 100
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		*Results
		Stats Statistics
	}{
		Results: results,
		Stats:   stats,
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
