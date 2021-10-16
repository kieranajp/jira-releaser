package handler

import (
	"net/url"

	"github.com/kieranajp/jira-releaser/pkg/github"
	"github.com/kieranajp/jira-releaser/pkg/jira"
	"github.com/urfave/cli/v2"
)

func Sync(c *cli.Context) error {
	u, err := url.ParseRequestURI(c.Args().Get(0))
	if err != nil {
		return cli.Exit("valid repo url must be provided", 1)
	}

	release, err := github.FetchRelease(u, c.String("release"))
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	tickets, err := github.ParseReleaseBody(release.Body)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	j, err := jira.New(c.String("jira-url"), c.String("jira-user"), c.String("jira-password"))
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return j.SetFixVersions(tickets, release)
}
