package cli

import (
	"fmt"
)

func (c *CLI) printUsage() {
	usage := `mycli is a load testing tool built with Go standard library.

Usage:
  mycli [command] [flags]

Available Commands:
  run         Run load test against target endpoint
  version     Print the version information
  help        Help about any command

Global Flags:
  --config string        config file (default is ./config.yaml)
  --rps int              requests per second (overrides config)
  --concurrency int      number of concurrent clients (overrides config)
  -o, --output string    output format: html, json (overrides config)

Use "mycli help [command]" for more information about a command.
`
	fmt.Fprint(c.stdout, usage)
}

func (c *CLI) helpCommand(command string) error {
	switch command {
	case "run":
		return c.runCommand([]string{"--help"})
	case "version":
		return c.versionCommand([]string{"--help"})
	default:
		fmt.Fprintf(c.stderr, "Unknown command: %s\n\n", command)
		c.printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}
