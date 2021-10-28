package main

import (
	"os"

	"github.com/kieranajp/jira-releaser/pkg/handler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Releaser",
		Usage: "Sync up Jira and Github releases",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "release",
				Aliases:  []string{"r"},
				Required: true,
				Usage:    "Github release version",
			},
			&cli.StringFlag{
				Name:    "github-user",
				Aliases: []string{"g"},
				Usage:   "Github username",
				EnvVars: []string{"GITHUB_USER"},
			},
			&cli.StringFlag{
				Name:    "github-token",
				Aliases: []string{"t"},
				Usage:   "Github token",
				EnvVars: []string{"GITHUB_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "jira-url",
				Aliases: []string{"j"},
				Value:   "https://jira.example.com",
				Usage:   "Jira URL",
				EnvVars: []string{"JIRA_URL"},
			},
			&cli.StringFlag{
				Name:    "jira-user",
				Aliases: []string{"u"},
				Usage:   "Jira Username",
				EnvVars: []string{"JIRA_USER"},
			},
			&cli.StringFlag{
				Name:    "jira-password",
				Aliases: []string{"p"},
				Usage:   "Jira Password",
				EnvVars: []string{"JIRA_PASS"},
			},
		},
		Action: handler.Sync,
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Exit")
	}
}
