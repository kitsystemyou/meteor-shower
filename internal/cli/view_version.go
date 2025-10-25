package cli

import (
	"flag"
	"fmt"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

func (c *CLI) versionCommand(args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	fs.Usage = func() {
		usage := `Display the version, git commit, and build date of the CLI tool.

Usage:
  meteor-shower version
`
		fmt.Fprint(c.stderr, usage)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Fprintf(c.stdout, "meteor-shower version %s\n", Version)
	fmt.Fprintf(c.stdout, "Git commit: %s\n", GitCommit)
	fmt.Fprintf(c.stdout, "Built: %s\n", BuildDate)

	return nil
}
