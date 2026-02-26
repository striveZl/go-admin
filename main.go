package main

import (
	"context"
	"fmt"
	"go-admin/cmd"
	"os"

	"github.com/urfave/cli/v3"
)

var VERSION = "1.0.0"

func main() {
	ctx := context.Background()

	app := &cli.Command{
		Name:    "goadmin",
		Version: VERSION,
		Usage:   "A Go admin project for beginners.",
		Commands: []*cli.Command{
			cmd.StartCmd(),
		},
	}

	err := app.Run(ctx, os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "application failed: %v\n", err)
		os.Exit(1)
	}
}
