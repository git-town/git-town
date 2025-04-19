# GitLab client-go (former `github.com/xanzy/go-gitlab`)

A GitLab API client enabling Go programs to interact with GitLab in a simple and uniform way.

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

For complete usage of go-gitlab, see the full [package docs](https://godoc.org/gitlab.com/gitlab-org/api/client-go).

## Testing

The `client-go` project comes with a `testing` package at `gitlab.com/gitlab-org/api/client-go/testing`
which contains a `TestClient` with [gomock](https://github.com/uber-go/mock) mocks for the individual services.

You can use them like this:

```go
func Test_MyApp(t *testing.T) {
    client := testing.NewTestClient(t)

    // Setup expectations
    client.MockClusterAgents.EXPECT().
        List(gomock.Any(), 123, nil).
        Return([]*gitlab.ClusterAgent{{ID: 1}}, nil)

    // Use the client in your test
    // You'd probably call your own code here that gets the client injected.
    // You can also retrieve a `gitlab.Client` object from `client.Client`.
    agents, err := client.ClusterAgents.List(ctx, 123, nil)
    assert.NoError(t, err)
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
