package jira

import (
	"regexp"

	j "github.com/andygrunwald/go-jira"
)

type JiraAPI struct {
	client *j.Client
}

func New(base, user, pass string) (*JiraAPI, error) {
	tp := j.BasicAuthTransport{
		Username: user,
		Password: pass,
	}

	client, err := j.NewClient(tp.Client(), base)
	if err != nil {
		return nil, err
	}

	return &JiraAPI{
		client: client,
	}, nil
}

func (c *JiraAPI) SetFixVersions(issueURLs []string, tagName string) error {
	for _, issueURL := range issueURLs {
		key := getKeyFromURL(issueURL)
		iss, err := c.getIssue(key)
		if err != nil {
			return err
		}

		iss.Fields.FixVersions = append(iss.Fields.FixVersions, &j.FixVersion{Name: tagName})
		_, _, err = c.client.Issue.Update(iss)
		if err != nil {
			return err
		}
	}

	return nil
}

func getKeyFromURL(issueURL string) string {
	r := regexp.MustCompile(`[\s|]?([A-Z]+-[0-9]+)[\s:|]?`)
	return r.FindString(issueURL)
}

func (c *JiraAPI) getIssue(key string) (*j.Issue, error) {
	issue, _, err := c.client.Issue.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}
