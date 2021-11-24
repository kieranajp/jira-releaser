package handler

import (
	"net/url"

	"github.com/kieranajp/jira-releaser/pkg/github"
	"github.com/kieranajp/jira-releaser/pkg/jira"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

func Sync(c *cli.Context) error {
	u, err := url.ParseRequestURI(c.Args().Get(0))
	if err != nil {
		return cli.Exit("valid repo url must be provided", 1)
	}

	pterm.DefaultSection.Printfln("Fetching %s from %s", c.String("release"), u.Path)

	gh := github.New(c.String("github-user"), c.String("github-token"))
	release, err := gh.FetchRelease(u, c.String("release"))
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	pterm.Println(release.Name)

	issues, err := github.ExtractIssues(release.Body)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	pterm.DefaultSection.WithLevel(2).Printfln("Found %d issues", len(issues))
	for _, issue := range issues {
		pterm.Println(issue)
	}

	j, err := jira.New(c.String("jira-url"), c.String("jira-user"), c.String("jira-password"))
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	pterm.DefaultSection.Println("Updating issues in JIRA")
	return j.SetFixVersions(issues, release)
}
