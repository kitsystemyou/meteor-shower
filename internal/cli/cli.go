package cli

import (
	"fmt"
	"io"
	"os"
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
