package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/example/mycli/internal/config"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

type CLI struct {
	args   []string
	stdout io.Writer
	stderr io.Writer
}

func New(args []string) *CLI {
	return &CLI{
		args:   args,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (c *CLI) Run() error {
	if len(c.args) < 1 {
		c.printUsage()
		return nil
	}

	command := c.args[0]

	switch command {
	case "run":
		return c.runCommand(c.args[1:])
	case "version":
		return c.versionCommand(c.args[1:])
	case "help":
		if len(c.args) > 1 {
			return c.helpCommand(c.args[1])
		}
		c.printUsage()
		return nil
	default:
		fmt.Fprintf(c.stderr, "Unknown command: %s\n\n", command)
		c.printUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

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

func (c *CLI) runCommand(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	configFile := fs.String("config", "", "config file (default is ./config.yaml)")
	output := fs.String("output", "text", "output format (text, json)")
	outputShort := fs.String("o", "text", "output format (text, json)")
	verbose := fs.Bool("verbose", false, "verbose output")
	verboseShort := fs.Bool("v", false, "verbose output")
	message := fs.String("message", "", "custom message to display")
	messageShort := fs.String("m", "", "custom message to display")

	fs.Usage = func() {
		usage := `Run executes the main application logic with the provided configuration.
You can specify a name as an argument, or it will use the name from the config file.

Usage:
  mycli run [name] [flags]

Flags:
  -m, --message string   custom message to display

Global Flags:
  --config string   config file (default is ./config.yaml)
  -o, --output string   output format (text, json) (default "text")
  -v, --verbose         verbose output
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

	name := cfg.App.Name
	if fs.NArg() > 0 {
		name = fs.Arg(0)
	}

	msg := ""
	if *message != "" {
		msg = *message
	} else if *messageShort != "" {
		msg = *messageShort
	} else if cfg.App.Message != "" {
		msg = cfg.App.Message
	} else {
		msg = "Hello"
	}

	outputFormat := *output
	if *outputShort != "text" {
		outputFormat = *outputShort
	}

	isVerbose := *verbose || *verboseShort

	result := map[string]interface{}{
		"name":    name,
		"message": fmt.Sprintf("%s, %s!", msg, name),
		"timeout": cfg.App.Timeout,
		"debug":   cfg.App.Debug,
	}

	if outputFormat == "json" {
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Fprintln(c.stdout, string(jsonData))
	} else {
		fmt.Fprintf(c.stdout, "%s\n", result["message"])
		if isVerbose {
			fmt.Fprintf(c.stdout, "Timeout: %ds\n", result["timeout"])
			fmt.Fprintf(c.stdout, "Debug: %v\n", result["debug"])
		}
	}

	return nil
}

func (c *CLI) versionCommand(args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	fs.Usage = func() {
		usage := `Display the version, git commit, and build date of the CLI tool.

Usage:
  mycli version
`
		fmt.Fprint(c.stderr, usage)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Fprintf(c.stdout, "mycli version %s\n", Version)
	fmt.Fprintf(c.stdout, "Git commit: %s\n", GitCommit)
	fmt.Fprintf(c.stdout, "Built: %s\n", BuildDate)

	return nil
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
