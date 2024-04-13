package giturl_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/shoenig/test/must"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := map[string]giturl.Parts{
		"git@github.com:git-town/git-town.git":                 {User: "git", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"git@bitbucket.org/git-town/git-town.git":              {User: "git", Host: "bitbucket.org", Org: "git-town", Repo: "git-town"},
		"git@bitbucket.org/git-town/git-town.github.com":       {User: "git", Host: "bitbucket.org", Org: "git-town", Repo: "git-town.github.com"},
		"git@github.com:git-town/git-town":                     {User: "git", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"git@gitlab.com:gitlab-com/www-gitlab-com.git":         {User: "git", Host: "gitlab.com", Org: "gitlab-com", Repo: "www-gitlab-com"},
		"git@gitlab.com:gitlab-com/www-gitlab-com":             {User: "git", Host: "gitlab.com", Org: "gitlab-com", Repo: "www-gitlab-com"},
		"git@gitlab.com:gitlab-org/quality/triage-ops.git":     {User: "git", Host: "gitlab.com", Org: "gitlab-org/quality", Repo: "triage-ops"},
		"git@gitlab.com:gitlab-org/quality/triage-ops":         {User: "git", Host: "gitlab.com", Org: "gitlab-org/quality", Repo: "triage-ops"},
		"username@bitbucket.org:git-town/git-town.git":         {User: "username", Host: "bitbucket.org", Org: "git-town", Repo: "git-town"},
		"username@bitbucket.org:git-town/git-town":             {User: "username", Host: "bitbucket.org", Org: "git-town", Repo: "git-town"},
		"https://github.com/git-town/git-town.git":             {User: "", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://github.com/git-town/git-town":                 {User: "", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user@github.com/git-town/git-town.git":        {User: "user", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user@github.com/git-town/git-town":            {User: "user", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user:secret@github.com/git-town/git-town.git": {User: "user:secret", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"https://user:secret@github.com/git-town/git-town":     {User: "user:secret", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://github.com/git-town/git-town.git":              {User: "", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://github.com/git-town/git-town":                  {User: "", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user@github.com/git-town/git-town.git":         {User: "user", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user@github.com/git-town/git-town":             {User: "user", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user:secret@github.com/git-town/git-town.git":  {User: "user:secret", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"http://user:secret@github.com/git-town/git-town":      {User: "user:secret", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"ssh://git@github.com/git-town/git-town.git":           {User: "git", Host: "github.com", Org: "git-town", Repo: "git-town"},
		"ssh://git@github.com/git-town/git-town":               {User: "git", Host: "github.com", Org: "git-town", Repo: "git-town"},
	}
	for give, want := range tests {
		have := giturl.Parse(give)
		must.EqOp(t, want, *have)
	}
}
