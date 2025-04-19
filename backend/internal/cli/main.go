package cli

import (
	"log"
	"nostalgia/internal/cli/actions"
	"os"

	"github.com/urfave/cli/v2"
)

func Main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "database",
				Aliases:  []string{"db"},
				Usage:    "Use sqlite database from `FILE`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "storage",
				Usage:    "Use `DIRECTORY` as storage directory for media",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "ffmpeg-binary",
				Usage:       "Use `PATH` as ffmpeg binary",
				DefaultText: "ffmpeg",
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "temporary-directory",
				Aliases:     []string{"temp-dir"},
				Usage:       "Use `PATH` as directory for temporary files",
				DefaultText: "./temp",
				Required:    false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "migrate-mediaviewer",
				Aliases: []string{"c"},
				Usage:   "migrate the old mediaviewer database to the new one",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "old-database",
						Aliases:  []string{"old-db"},
						Usage:    "Use old database from `FILE`",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "old-storage",
						Usage:    "Use old storage from `DIRECTORY`",
						Required: true,
					},
				},
				Action: actions.MigrateMediaviewer,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
