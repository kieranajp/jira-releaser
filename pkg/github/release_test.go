package github_test

import (
	"testing"

	g "github.com/kieranajp/jira-releaser/pkg/github"
)

func TestParsingLinks(t *testing.T) {
	body := "Released the thing\n\n\n- [JIRA-123](https://jira.com/123)\n- Fixed [JIRA-234](https://jira.com/234)\n\nAnd other things"
	expected := []string{"https://jira.com/123", "https://jira.com/234"}

	actual, err := g.ParseReleaseBody(body)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(actual) != len(expected) {
		t.Errorf("Expected %d issues, got %d", len(expected), len(actual))
	}

	if actual[0] != expected[0] {
		t.Errorf("Expected %s, got %s", expected[0], actual[0])
	}

	if actual[1] != expected[1] {
		t.Errorf("Expected %s, got %s", expected[1], actual[1])
	}
}
