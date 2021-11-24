package jira

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	j "github.com/andygrunwald/go-jira"
	"github.com/kieranajp/jira-releaser/pkg/github"
	"github.com/pterm/pterm"
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

func (c *JiraAPI) SetFixVersions(issues []string, release *github.Release) error {
	p, _ := pterm.DefaultProgressbar.WithTotal(len(issues)).WithTitle("Releasing issues").Start()
	failed := make([]string, 0)

	for _, key := range issues {
		p.UpdateTitle(fmt.Sprintf("Fetching issue %s", key))
		iss, err := c.getIssue(key)
		if err != nil {
			pterm.Warning.Printfln("Unable to find issue %s", key)
			failed = append(failed, key)
			continue
		}

		time.Sleep(time.Millisecond * 300)

		version, err := c.ensureFixVersionExists(iss, release)
		if err != nil {
			return err
		}

		p.UpdateTitle(fmt.Sprintf("Setting FixVersion for %s", key))
		err = c.addVersionToIssue(iss, version)
		if err != nil {
			failed = append(failed, key)
			pterm.Warning.Printfln("Unable to set fix version for %s", key)
			continue
		}
		pterm.Success.Printfln("Updated %s", key)

		p.Increment()
	}

	p.Stop()

	printReport(issues, failed)

	return nil
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

	_, repo := parseRepoNameFromURL(release.URL)
	fixName := fmt.Sprintf("%s %s", repo, release.TagName)
	version := &j.Version{
		ProjectID:   projectID,
		Name:        fixName,
		Description: fmt.Sprintf("%s\n\n(%s)", release.Name, release.URL),
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

func printReport(issues, failed []string) {
	pterm.DefaultSection.WithLevel(2).Println("Issue report")
	d := pterm.TableData{{"Issue", "Status"}}
	for _, s := range issues {
		if stringInSlice(s, failed) {
			d = append(d, []string{pterm.LightRed(s), pterm.LightRed("failed")})
		} else {
			d = append(d, []string{s, pterm.LightGreen("updated")})
		}
	}
	pterm.DefaultTable.WithHasHeader().WithData(d).Render()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
