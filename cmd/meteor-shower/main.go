package main

import (
	"os"

	"github.com/kitsystemyou/meteor-shower/internal/cli"
)

func main() {
	app := cli.New(os.Args[1:])
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
