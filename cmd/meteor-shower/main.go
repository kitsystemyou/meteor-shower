package main

import (
	"fmt"
	"os"

	"github.com/kitsystemyou/meteor-shower/internal/cli"
)

func main() {
	app := cli.New(os.Args[1:])
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
