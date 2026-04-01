package cmd

import (
	"context"
	"strings"

	"github.com/urfave/cli/v3"

	"go-admin/internal/config"
	"go-admin/internal/db"
)

func MigrateCmd() *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "run database migrations",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "workdir",
				Aliases:     []string{"d"},
				Usage:       "Working directory",
				DefaultText: "configs",
				Value:       "configs",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Runtime configuration files or directory (relative to workdir, multiple separated by commas)",
				DefaultText: "dev",
				Value:       "dev",
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			workdir := command.String("workdir")
			configNames := parseConfigArg(command.String("config"))
			config.MustLoad(workdir, configNames...)
			config.C.General.WorkDir = workdir

			dbStore, err := db.Open(ctx, config.C.Database)
			if err != nil {
				return err
			}
			defer func() {
				_ = dbStore.Close()
			}()

			return db.Migrate(ctx, dbStore.SQL())
		},
	}
}

func parseConfigArg(raw string) []string {
	items := strings.Split(raw, ",")
	names := make([]string, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}
