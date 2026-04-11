# GitLab client-go (former `github.com/xanzy/go-gitlab`)

A GitLab API client enabling Go programs to interact with GitLab in a simple and uniform way.


## Table of Contents

[[_TOC_]]

## Usage

```go
import "gitlab.com/gitlab-org/api/client-go"
```

Construct a new GitLab client, then use the various services on the client to
access different parts of the GitLab API. For example, to list all
users:

```go
git, err := gitlab.NewClient("yourtokengoeshere")
if err != nil {
  log.Fatalf("Failed to create client: %v", err)
}
users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{})
```

There are a few `With...` option functions that can be used to customize
the API client. For example, to set a custom base URL:

```go
git, err := gitlab.NewClient("yourtokengoeshere", gitlab.WithBaseURL("https://git.mydomain.com/api/v4"))
if err != nil {
  log.Fatalf("Failed to create client: %v", err)
}
users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{})
```

Some API methods have optional parameters that can be passed. For example,
to list all projects for user "svanharmelen":

```go
git := gitlab.NewClient("yourtokengoeshere")
opt := &gitlab.ListProjectsOptions{Search: gitlab.Ptr("svanharmelen")}
projects, _, err := git.Projects.ListProjects(opt)
```

### Examples

The [examples](/examples) directory
contains a couple for clear examples, of which one is partially listed here as well:

```go
package main

import (
	"log"

	"gitlab.com/gitlab-org/api/client-go"
)

func main() {
	git, err := gitlab.NewClient("yourtokengoeshere")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create new project
	p := &gitlab.CreateProjectOptions{
		Name:                     gitlab.Ptr("My Project"),
		Description:              gitlab.Ptr("Just a test project to play with"),
		MergeRequestsAccessLevel: gitlab.Ptr(gitlab.EnabledAccessControl),
		SnippetsAccessLevel:      gitlab.Ptr(gitlab.EnabledAccessControl),
		Visibility:               gitlab.Ptr(gitlab.PublicVisibility),
	}
	project, _, err := git.Projects.CreateProject(p)
	if err != nil {
		log.Fatal(err)
	}

	// Add a new snippet
	s := &gitlab.CreateProjectSnippetOptions{
		Title:           gitlab.Ptr("Dummy Snippet"),
		FileName:        gitlab.Ptr("snippet.go"),
		Content:         gitlab.Ptr("package main...."),
		Visibility:      gitlab.Ptr(gitlab.PublicVisibility),
	}
	_, _, err = git.ProjectSnippets.CreateSnippet(project.ID, s)
	if err != nil {
		log.Fatal(err)
	}
}
```

#### Use OAuth2 helper package

The following example demonstrates how to use the `gitlab.com/gitlab-org/api/client-go/oauth2` package:

```go
package main

import (
	"context"
	"fmt"
	"os/exec"

	gitlab "gitlab.com/gitlab-org/api/client-go"
	"gitlab.com/gitlab-org/api/client-go/gitlaboauth2"
)

func main() {
	ctx := context.Background()
	// Authorize with GitLab.com and OAuth2
	clientID := "aaa"
	redirectURL := "http://localhost:9999/auth/redirect"
	scopes := []string{"read_api"}
	config := gitlaboauth2.NewOAuth2Config("", clientID, redirectURL, scopes)

	server := gitlaboauth2.NewCallbackServer(config, ":9999", func(url string) error {
		return exec.Command("open", url).Start()
	})

	token, err := server.GetToken(ctx)
	if err != nil {
		panic(err)
	}

	client, err := gitlab.NewAuthSourceClient(gitlab.OAuthTokenSource{TokenSource: config.TokenSource(ctx, token)})
	if err != nil {
		panic(err)
	}

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current user: %s\n", user.Username)
}
```

#### Use the `config` package (experimental)

The `config` package defines a configuration file format (YAML) to configure GitLab instances
and their associated authentication methods combined in contexts (similar to what you might know from Kubernetes).

The configuration is located in the users config directory (e.g. XDG config dir), in `gitlab/config.yaml`.

A basic example for an OAuth flow for GitLab.com that stores the credentials in the systems keyring, looks like this:

```yaml
version: gitlab.com/config/v1beta1

instances:
    - name: gitlab-com
      server: https://gitlab.com

auths:
    - name: oauth-keyring
      auth-info:
        oauth2:
            access-token-source:
                keyring:
                    service: client-go
                    user: access-token
            refresh-token-source:
                keyring:
                    service: client-go
                    user: refresh-token
contexts:
    - name: gitlab-com-keyring
      instance: gitlab-com
      auth: oauth-keyring

current-context: gitlab-com-keyring
```

An application with `client-go` is able to effortlessly create a new client using that configuration:

```go
package main

import (
	"fmt"
	"log"

	"gitlab.com/gitlab-org/api/client-go"
	"gitlab.com/gitlab-org/api/client-go/config"
)

func main() {
	// Create a config with default location (~/.config/gitlab/config.yaml)
	cfg := config.New(
		config.WithOAuth2Settings(config.OAuth2Settings{
			AuthorizationFlowEnabled: true,
			CallbackServerListenAddr: ":7171",
			Browser: func(url string) error {
				fmt.Printf("Open: %s\n", url)
				return nil
			},
			ClientID:    "<your-client-id>",
			RedirectURL: "http://localhost:7171/auth/redirect",
			Scopes:      []string{"read_api"},
		}),
	)

	// Load the configuration
	if err := cfg.Load(); err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	client, err := cfg.NewClient(gitlab.WithUserAgent("my-app"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use the client
	user, _, err := client.Users.CurrentUser()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

	fmt.Printf("Authenticated as: %s (%s)\n", user.Name, user.Username)
}
```

For complete usage of go-gitlab, see the full [package docs](https://godoc.org/gitlab.com/gitlab-org/api/client-go).

## Installation

To install the library, use the following command:

```go
go get gitlab.com/gitlab-org/api/client-go
```

## Testing

The `client-go` project comes with a `testing` package at `gitlab.com/gitlab-org/api/client-go/testing`
which contains a `TestClient` with [gomock](https://github.com/uber-go/mock) mocks for the individual services.

You can use them like this:

```go
func TestMockExample(t *testing.T) {
    client := gitlabtesting.NewTestClient(t)
    opts := &gitlab.ListAgentsOptions{}
    expectedResp := &gitlab.Response{}
    pid := 1
    // Setup expectations
    client.MockClusterAgents.EXPECT().
        ListAgents(pid, opts).
        Return([]*gitlab.Agent{{ID: 1}}, expectedResp, nil)

    // Use the client in your test
    // You'd probably call your own code here that gets the client injected.
    // You can also retrieve a `gitlab.Client` object from `client.Client`.
    agents, resp, err := client.ClusterAgents.ListAgents(pid, opts)
    assert.NoError(t, err)
    assert.Equal(t, expectedResp, resp)
    assert.Len(t, agents, 1)
}
```

### I want to generate my own mocks

You can! You can set up your own `TestClient` with mocks pretty easily:

```go
func NewTestClient(t *testing.T) {
    // generate your mocks or instantiate a fake or whatever you like
    mockClusterAgentsService := newMockClusterAgentsService(t)
	client := &gitlab.Client{
		ClusterAgents: mockClusterAgentsService
	}

	return tc
}
```

The `newMockClusterAgentsService` must return a type that implements `gitlab.ClusterAgentsInterface`.

You can have a look at [`testing/client.go`](/testing.client.go) how it's implemented for `gomock`.

## Compatibility

The `client-go` package will maintain compatibility with the officially supported Go releases
at the time the package is released. According to the [Go Release Policy](https://go.dev/doc/devel/release#policy),
that's currently the two last major Go releases.
This compatibility is reflected in the `go` directive of the [`go.mod`](/go.mod) file
and the unit test matrix in [`.gitlab-ci.yml`](/.gitlab-ci.yml).

You may also use https://endoflife.date/go to quickly discover the supported Go versions.

## Contributing

Contributions are always welcome. For more information, check out the
[contributing guide](/CONTRIBUTING.md).

## Maintenance

This is a community maintained project. If you have a paid GitLab subscription,
please note that this project is not packaged as a part of GitLab, and falls outside
of the scope of support.

For more information, see GitLab's
[Statement of Support](https://about.gitlab.com/support/statement-of-support.html).
Please fill out an issue in this projects issue tracker and someone from the community
will respond as soon as they are available to help you.

### Known GitLab Projects using this package

- [GitLab Terraform Provider](https://gitlab.com/gitlab-org/terraform-provider-gitlab)
  maintained by the community with support from ~"group::environments"
- [GitLab CLI (`glab`)](https://gitlab.com/gitlab-org/cli)
  maintained by ~"group::code review"
