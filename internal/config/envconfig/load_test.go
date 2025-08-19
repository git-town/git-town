package envconfig_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestLoad(t *testing.T) {
	tokenText := "my-github-token"
	authTokenText := "my-github-auth-token"

	tests := []struct {
		name        string
		githubToken Option[string]
		authToken   Option[string]
		want        Option[forgedomain.GitHubToken]
	}{
		{
			name:        "loads from GITHUB_TOKEN when both are set",
			githubToken: Some(tokenText),
			authToken:   Some(authTokenText),
			want:        Some(forgedomain.GitHubToken(tokenText)),
		},
		{
			name:        "loads from GITHUB_AUTH_TOKEN if GITHUB_TOKEN is empty",
			githubToken: None[string](),
			authToken:   Some(authTokenText),
			want:        Some(forgedomain.GitHubToken(authTokenText)),
		},
		{
			name:        "loads from GITHUB_TOKEN only",
			githubToken: Some(tokenText),
			authToken:   None[string](),
			want:        Some(forgedomain.GitHubToken(tokenText)),
		},
		{
			name:        "returns none when no env is set",
			githubToken: None[string](),
			authToken:   None[string](),
			want:        None[forgedomain.GitHubToken](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv("GITHUB_TOKEN")
			os.Unsetenv("GITHUB_AUTH_TOKEN")
			if githubToken, has := tt.githubToken.Get(); has {
				t.Setenv("GITHUB_TOKEN", githubToken)
			}
			if authToken, has := tt.authToken.Get(); has {
				t.Setenv("GITHUB_AUTH_TOKEN", authToken)
			}
			cfg := envconfig.Load()
			must.Eq(t, tt.want, cfg.GitHubToken)
		})
	}
}
