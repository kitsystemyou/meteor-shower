package cli

import (
	"flag"
	"fmt"
	"os"
)

const defaultConfigTemplate = `# meteor-shower configuration file for load testing
loadtest:
  # Target domain
  domain: "http://localhost:8080"
  
  # Endpoints with weights
  # Weight determines the distribution of requests across endpoints
  endpoints:
    - path: "/"
      weight: 1.0
  
  # Example: Multiple endpoints with different weights
  # endpoints:
  #   - path: "/"
  #     weight: 1.0      # Highest frequency
  #   - path: "/health"
  #     weight: 0.5      # Medium frequency
  #   - path: "/slow"
  #     weight: 0.2      # Lowest frequency
  
  # Requests per second
  rps: 10
  
  # Number of concurrent clients
  concurrency: 1
  
  # Test duration in seconds
  duration: 10
  
  # Output format (html or json)
  output: "html"
`

func (c *CLI) configInitCommand(args []string) error {
	fs := flag.NewFlagSet("config init", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	outputFile := fs.String("output", "config.yaml", "output file path")
	outputShort := fs.String("o", "config.yaml", "output file path")
	force := fs.Bool("force", false, "overwrite existing file")
	forceShort := fs.Bool("f", false, "overwrite existing file")

	fs.Usage = func() {
		usage := `Generate a default configuration file.

Usage:
  meteor-shower config init [flags]

Flags:
  -o, --output string    output file path (default "config.yaml")
  -f, --force            overwrite existing file

Examples:
  # Generate config.yaml in current directory
  meteor-shower config init

  # Generate with custom filename
  meteor-shower config init -o my-config.yaml

  # Overwrite existing file
  meteor-shower config init -f
`
		fmt.Fprint(c.stderr, usage)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Determine output file
	outFile := *outputFile
	if *outputShort != "config.yaml" {
		outFile = *outputShort
	}

	// Determine force flag
	forceFlag := *force || *forceShort

	// Check if file exists
	if _, err := os.Stat(outFile); err == nil && !forceFlag {
		return fmt.Errorf("file %s already exists. Use -f to overwrite", outFile)
	}

	// Write config file
	if err := os.WriteFile(outFile, []byte(defaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Fprintf(c.stdout, "Configuration file created: %s\n", outFile)
	return nil
}

func (c *CLI) configCommand(args []string) error {
	if len(args) == 0 {
		fmt.Fprintf(c.stderr, "Error: subcommand required\n\n")
		c.printConfigUsage()
		return fmt.Errorf("subcommand required")
	}

	subcommand := args[0]

	switch subcommand {
	case "init":
		return c.configInitCommand(args[1:])
	default:
		fmt.Fprintf(c.stderr, "Error: unknown subcommand: %s\n\n", subcommand)
		c.printConfigUsage()
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

func (c *CLI) printConfigUsage() {
	usage := `Manage configuration files.

Usage:
  meteor-shower config [subcommand] [flags]

Available Subcommands:
  init        Generate a default configuration file

Use "meteor-shower config [subcommand] --help" for more information about a subcommand.
`
	fmt.Fprint(c.stderr, usage)
}
