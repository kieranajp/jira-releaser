package jira

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	j "github.com/andygrunwald/go-jira"
	"github.com/kieranajp/jira-releaser/pkg/github"
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

func (c *JiraAPI) SetFixVersions(issueURLs []string, release *github.Release) error {
	for _, issueURL := range issueURLs {
		key := getKeyFromURL(issueURL)
		iss, err := c.getIssue(key)
		if err != nil {
			return err
		}

		version, err := c.ensureFixVersionExists(iss, release)
		if err != nil {
			return err
		}

		err = c.addVersionToIssue(iss, version)
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

func parseRepoNameFromURL(url string) (string, string) {
	r := regexp.MustCompile(`https:\/\/.+?\/(.+?)\/(.+?)\/`)
	matches := r.FindStringSubmatch(url)
	return matches[1], matches[2]
}

func (c *JiraAPI) ensureFixVersionExists(iss *j.Issue, release *github.Release) (*j.Version, error) {
	projectID, _ := strconv.Atoi(iss.Fields.Project.ID)

	owner, repo := parseRepoNameFromURL(release.URL)
	fixName := fmt.Sprintf("%s/%s %s", owner, repo, release.TagName)
	version := &j.Version{
		ProjectID:   projectID,
		Name:        fixName,
		Description: fmt.Sprintf("%s\n\n(%s)", release.Body, release.URL),
		ReleaseDate: release.PublishedAt,
	}

	_, resp, err := c.client.Version.Create(version)
	if resp.StatusCode == 400 {
		return version, nil
	}
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (c *JiraAPI) addVersionToIssue(issue *j.Issue, version *j.Version) error {
	reader := strings.NewReader(fmt.Sprintf(
		`{"update": {"fixVersions": [{"add": {"name": "%s"}}]}}`,
		version.Name))
	req, err := c.client.NewRawRequest("PUT", fmt.Sprintf("rest/api/2/issue/%s", issue.ID), reader)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
