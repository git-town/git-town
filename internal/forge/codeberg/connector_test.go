package codeberg_test

import (
	"testing"

	forgejosdk "codeberg.org/mvdkleijn/forgejo-sdk"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/forge/codeberg"
	"github.com/shoenig/test/must"
)

func TestFilterCodebergPullRequests(t *testing.T) {
	t.Parallel()
	give := []*forgejosdk.PullRequest{
		// matching branch
		{
			Head: &forgejosdk.PRBranchInfo{
				Name: "branch",
			},
			Base: &forgejosdk.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different name
		{
			Head: &forgejosdk.PRBranchInfo{
				Name: "other",
			},
			Base: &forgejosdk.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different target
		{
			Head: &forgejosdk.PRBranchInfo{
				Name: "branch",
			},
			Base: &forgejosdk.PRBranchInfo{
				Name: "other",
			},
		},
	}
	want := []*forgejosdk.PullRequest{
		{
			Head: &forgejosdk.PRBranchInfo{
				Name: "branch",
			},
			Base: &forgejosdk.PRBranchInfo{
				Name: "target",
			},
		},
	}
	have := codeberg.FilterPullRequests(give, "branch", "target")
	must.Eq(t, want, have)
}

//nolint:paralleltest  // mocks HTTP
func TestCodeberg(t *testing.T) {
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		give := forgedomain.Proposal{
			Number: 1,
			Title:  "my title",
		}
		want := "my title (#1)"
		connector := codeberg.Connector{}
		have := connector.DefaultProposalMessage(give)
		must.EqOp(t, want, have)
	})

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE CODEBERG CONNECTOR.
	//
	// t.Run("NewProposalURL", func(t *testing.T) {
	// 	connector, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingCodeberg,
	// 		RemoteURL:      giturl.Parse("git@codeberg.org:git-town/docs.git"),
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.NoError(t, err)
	// 	have, err := connector.NewProposalURL("feature", "parent")
	// 	must.NoError(t, err)
	// 	must.EqOp(t, "https://codeberg.org/git-town/docs/compare/parent...feature", have)
	// })

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE CODEBERG CONNECTOR.
	//
	// t.Run("RepositoryURL", func(t *testing.T) {
	// 	connector, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingCodeberg,
	// 		RemoteURL:      giturl.Parse("git@codeberg.org:git-town/docs.git"),
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.NoError(t, err)
	// 	have := connector.RepositoryURL()
	// 	must.EqOp(t, "https://codeberg.org/git-town/docs", have)
	// })
}

func TestNewCodebergConnector(t *testing.T) {
	t.Parallel()

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE CODEBERG CONNECTOR.
	//
	// t.Run("hosted service type provided manually", func(t *testing.T) {
	// 	t.Parallel()
	// 	have, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingCodeberg,
	// 		RemoteURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
	// 		APIToken:       "apiToken",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.NoError(t, err)
	// 	wantConfig := hostingdomain.Config{
	// 		Hostname:     "custom-url.com",
	// 		Organization: "git-town",
	// 		Repository:   "docs",
	// 	}
	// 	must.EqOp(t, wantConfig, have.Config)
	// })

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE CODEBERG CONNECTOR.
	//
	// t.Run("repo is hosted by another forge type --> no connector", func(t *testing.T) {
	// 	t.Parallel()
	// 	have, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingNone,
	// 		RemoteURL:      giturl.Parse("git@github.com:git-town/git-town.git"),
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.Nil(t, have)
	// 	must.NoError(t, err)
	// })

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE CODEBERG CONNECTOR.
	//
	// t.Run("no origin remote --> no connector", func(t *testing.T) {
	// 	t.Parallel()
	// 	var remoteURL *giturl.Parts
	// 	have, err := codeberg.NewConnector(codeberg.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingNone,
	// 		RemoteURL:      remoteURL,
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.Nil(t, have)
	// 	must.NoError(t, err)
	// })
}
