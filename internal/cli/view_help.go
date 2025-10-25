package cli

import (
	"fmt"
)

func (c *CLI) printUsage() {
	usage := `mycli is a sample CLI tool that demonstrates:
- Command-line argument parsing
- YAML configuration file support
- Subcommands with help documentation
- Best practices for Go CLI applications

Usage:
  mycli [command] [flags]

Available Commands:
  run         Run the main application logic
  version     Print the version information
  help        Help about any command

Global Flags:
  --config string   config file (default is ./config.yaml)
  -o, --output string   output format (text, json) (default "text")
  -v, --verbose         verbose output

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
