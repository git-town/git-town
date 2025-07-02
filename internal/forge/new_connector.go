package forge

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v21/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v21/internal/forge/codeberg"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/gh"
	"github.com/git-town/git-town/v21/internal/forge/gitea"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/forge/gitlab"
	"github.com/git-town/git-town/v21/internal/forge/glab"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// NewConnector provides an instance of the forge connector to use based on the given gitConfig.
func NewConnector(args NewConnectorArgs) (Option[forgedomain.Connector], error) {
	remoteURL, hasRemoteURL := args.RemoteURL.Get()
	forgeType, hasForgeType := Detect(remoteURL, args.ForgeType).Get()
	if !hasRemoteURL || !hasForgeType {
		return None[forgedomain.Connector](), nil
	}
	var connector forgedomain.Connector
	var err error
	switch forgeType {
	case forgedomain.ForgeTypeBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			AppPassword: args.BitbucketAppPassword,
			ForgeType:   args.ForgeType,
			Log:         args.Log,
			RemoteURL:   remoteURL,
			UserName:    args.BitbucketUsername,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeBitbucketDatacenter:
		connector = bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			AppPassword: args.BitbucketAppPassword,
			ForgeType:   args.ForgeType,
			Log:         args.Log,
			RemoteURL:   remoteURL,
			UserName:    args.BitbucketUsername,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeCodeberg:
		connector, err = codeberg.NewConnector(codeberg.NewConnectorArgs{
			APIToken:  args.CodebergToken,
			Log:       args.Log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	case forgedomain.ForgeTypeGitea:
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  args.GiteaToken,
			Log:       args.Log,
			RemoteURL: remoteURL,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeGitHub:
		if githubConnectorType, hasGitHubConnectorType := args.GitHubConnectorType.Get(); hasGitHubConnectorType {
			switch githubConnectorType {
			case forgedomain.GitHubConnectorTypeAPI:
				connector, err = github.NewConnector(github.NewConnectorArgs{
					APIToken:  args.GitHubToken,
					Log:       args.Log,
					RemoteURL: remoteURL,
				})
				return Some(connector), err
			case forgedomain.GitHubConnectorTypeGh:
				connector = gh.Connector{
					Backend:  args.Backend,
					Frontend: args.Frontend,
				}
				return Some(connector), err
			}
		}
		// no GitHubConnectorType specified --> use the API connector
		connector, err = github.NewConnector(github.NewConnectorArgs{
			APIToken:  args.GitHubToken,
			Log:       args.Log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	case forgedomain.ForgeTypeGitLab:
		if gitLabConnectorType, hasGitLabConnectorType := args.GitLabConnectorType.Get(); hasGitLabConnectorType {
			switch gitLabConnectorType {
			case forgedomain.GitLabConnectorTypeAPI:
				connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
					APIToken:  args.GitLabToken,
					Log:       args.Log,
					RemoteURL: remoteURL,
				})
				return Some(connector), err
			case forgedomain.GitLabConnectorTypeGlab:
				connector = glab.Connector{
					Backend:  args.Backend,
					Frontend: args.Frontend,
				}
				return Some(connector), err
			}
		}
		// no GitLabConnectorType specified --> use the API connector
		connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  args.GitLabToken,
			Log:       args.Log,
			RemoteURL: remoteURL,
		})
		return Some(connector), err
	}
	return None[forgedomain.Connector](), nil
}

type NewConnectorArgs struct {
	Backend              subshelldomain.Querier
	BitbucketAppPassword Option[forgedomain.BitbucketAppPassword]
	BitbucketUsername    Option[forgedomain.BitbucketUsername]
	CodebergToken        Option[forgedomain.CodebergToken]
	ForgeType            Option[forgedomain.ForgeType]
	Frontend             subshelldomain.Runner
	GitHubConnectorType  Option[forgedomain.GitHubConnectorType]
	GitHubToken          Option[forgedomain.GitHubToken]
	GitLabConnectorType  Option[forgedomain.GitLabConnectorType]
	GitLabToken          Option[forgedomain.GitLabToken]
	GiteaToken           Option[forgedomain.GiteaToken]
	Log                  print.Logger
	RemoteURL            Option[giturl.Parts]
}
