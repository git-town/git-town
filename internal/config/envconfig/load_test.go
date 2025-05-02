package envconfig_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/config/envconfig"
	"github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLoad(t *testing.T) {
	tokenText := "my-github-token"
	authTokenText := "my-github-auth-token"

	tests := []struct {
		name        string
		githubToken *string
		authToken   *string
		want        prelude.Option[configdomain.GitHubToken]
	}{
		{
			name:        "loads from GITHUB_TOKEN when both are set",
			githubToken: &tokenText,
			authToken:   &authTokenText,
			want:        prelude.Some(configdomain.GitHubToken(tokenText)),
		},
		{
			name:        "loads from GITHUB_AUTH_TOKEN if GITHUB_TOKEN is empty",
			githubToken: ptr(""),
			authToken:   &authTokenText,
			want:        prelude.Some(configdomain.GitHubToken(authTokenText)),
		},
		{
			name:        "loads from GITHUB_TOKEN only",
			githubToken: &tokenText,
			authToken:   nil,
			want:        prelude.Some(configdomain.GitHubToken(tokenText)),
		},
		{
			name:        "returns none when no env is set",
			githubToken: nil,
			authToken:   nil,
			want:        prelude.None[configdomain.GitHubToken](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv("GITHUB_TOKEN")
			os.Unsetenv("GITHUB_AUTH_TOKEN")

			if tt.githubToken != nil {
				t.Setenv("GITHUB_TOKEN", *tt.githubToken)
			}
			if tt.authToken != nil {
				t.Setenv("GITHUB_AUTH_TOKEN", *tt.authToken)
			}

			cfg := envconfig.Load()
			must.Eq(t, tt.want, cfg.GitHubToken)
		})
	}
}

func ptr(s string) *string {
	return &s
}
