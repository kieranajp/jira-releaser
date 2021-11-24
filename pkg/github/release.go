package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type User struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	URL       string `json:"html_url"`
}

type Release struct {
	URL         string `json:"html_url,omitempty"`
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Body        string `json:"body,omitempty"`
	TagName     string `json:"tag_name,omitempty"`
	Author      *User  `json:"author,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
}

type Github struct {
	user, token string
}

func New(user, token string) *Github {
	return &Github{
		user:  user,
		token: token,
	}
}

func (g *Github) FetchRelease(repo *url.URL, tag string) (*Release, error) {
	var release Release
	url := fmt.Sprintf("https://api.github.com/repos%s/releases/tags/%s", repo.Path, tag)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(g.user, g.token)
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP response from Github: %s", resp.Status)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func ExtractIssues(body string) ([]string, error) {
	re := regexp.MustCompile(`\b[A-Z]{2,}-\d+\b`)

	matches := re.FindAllString(body, -1)
	if len(matches) == 0 {
		return nil, errors.New("no issues found in release body")
	}

	return matches, nil
}
