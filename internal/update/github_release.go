package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	GitHubAPIURL = "https://api.github.com/repos/git-town/git-town/releases/latest"
	HTTPTimeout  = 5 * time.Second
)

// GitHubRelease represents the structure of a GitHub release API response
type GitHubRelease struct {
	HTMLURL string `json:"html_url"`
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
}

// GitHubClient handles interactions with the GitHub Releases API
type GitHubClient struct {
	httpClient *http.Client
}

// NewGitHubClient creates a new GitHub client with a timeout
func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{
			Timeout: HTTPTimeout,
		},
	}
}

// GetLatestRelease fetches the latest release from GitHub
func (self *GitHubClient) GetLatestRelease(ctx context.Context) (*GitHubRelease, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, GitHubAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "git-town")

	resp, err := self.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &release, nil
}
