package hosting

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v7/src/config"
)

// BitbucketConnector provides access to the API of Bitbucket installations.
type AzureDevopsConnector struct {
	CommonConfig
	MainBranch string
}

// NewBitbucketConnector provides a Bitbucket connector instance if the current repo is hosted on Bitbucket,
// otherwise nil.
func NewAzureDevopsConnector(gitConfig gitTownConfig) (*AzureDevopsConnector, error) {
	hostingService, err := gitConfig.HostingService()
	if err != nil {
		return nil, err
	}
	url := gitConfig.OriginURL()
	fmt.Println("11111111111", url)
	fmt.Println("22222222222", url.Host)
	if url == nil || (url.Host != "dev.azure.com" && hostingService != config.HostingServiceAzureDevops) {
		return nil, nil //nolint:nilnil
	}
	return &AzureDevopsConnector{
		CommonConfig: CommonConfig{
			APIToken:     "",
			Hostname:     url.Host,
			Organization: url.Org,
			Repository:   url.Repo,
		},
	}, nil
}

func (c *AzureDevopsConnector) FindProposal(branch, target string) (*Proposal, error) {
	return nil, fmt.Errorf("the Azure Devops API functionality isn't implemented yet")
}

func (c *AzureDevopsConnector) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (#%d)", proposal.Title, proposal.Number)
}

func (c *AzureDevopsConnector) HostingServiceName() string {
	return "Azure DevOps"
}

func (c *AzureDevopsConnector) NewProposalURL(branch, parentBranch string) (string, error) {
	toCompare := branch
	if parentBranch != c.MainBranch {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("%s/compare/%s?expand=1", c.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (c *AzureDevopsConnector) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", c.Hostname, c.Organization, c.Repository)
}

//nolint:nonamedreturns
func (c *AzureDevopsConnector) SquashMergeProposal(number int, message string) (mergeSHA string, err error) {
	return "", errors.New("shipping pull requests via the Azure DevOps API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (c *AzureDevopsConnector) UpdateProposalTarget(number int, target string) error {
	return errors.New("shipping pull requests via the Azure DevOps API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}
