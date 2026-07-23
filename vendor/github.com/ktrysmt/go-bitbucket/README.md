# go-bitbucket

<a class="repo-badge" href="https://godoc.org/github.com/ktrysmt/go-bitbucket"><img src="https://godoc.org/github.com/ktrysmt/go-bitbucket?status.svg" alt="go-bitbucket?status"></a>
<a href="https://goreportcard.com/report/github.com/ktrysmt/go-bitbucket"><img class="badge" tag="github.com/ktrysmt/go-bitbucket" src="https://goreportcard.com/badge/github.com/ktrysmt/go-bitbucket"></a>

> Bitbucket-API library for golang.

Supports Bitbucket Cloud REST API v2.0. Responses follow the JSON shape
documented by the official API.

- Bitbucket Cloud REST API v2.0 <https://developer.atlassian.com/cloud/bitbucket/rest/intro/>
- Swagger for API v2.0 <https://api.bitbucket.org/swagger.json>

## Install

```sh
go get github.com/ktrysmt/go-bitbucket
```

## Authentication

Bitbucket Cloud accepts several credential types. Pick the constructor that
matches your credential.

### API token (recommended)

Atlassian has deprecated app passwords; use an Atlassian API token instead.
Pass your Atlassian account email as the username and the API token as the
password.

```go
c := bitbucket.NewAPITokenAuth("you@example.com", "your-api-token")
```

`NewAPITokenAuth` is a thin alias over `NewBasicAuth` that documents the
intended usage. See the
[Atlassian API token guide](https://support.atlassian.com/bitbucket-cloud/docs/using-api-tokens/).

### App password (legacy)

```go
c := bitbucket.NewBasicAuth("username", "app-password")
```

### OAuth bearer token

```go
c, err := bitbucket.NewOAuthbearerToken("access-token")
```

### OAuth flow (client credentials, authorization code, refresh token)

If you obtain tokens through the OAuth handshake yourself, use one of:

```go
c, err := bitbucket.NewOAuthClientCredentials("client-id", "client-secret")
c, accessToken, err := bitbucket.NewOAuthWithCode("client-id", "client-secret", "auth-code")
c, accessToken, err := bitbucket.NewOAuthWithRefreshToken("client-id", "client-secret", "refresh-token")
```

For Isolated Cloud Instances whose OAuth token endpoint lives under a
customer-specific hostname, the matching `*WithEndpoint` variants override the
default `bitbucket.Endpoint`:

```go
import "golang.org/x/oauth2"

ep := oauth2.Endpoint{
    AuthURL:  "https://auth.your-isolated-instance.example.com/site/oauth2/authorize",
    TokenURL: "https://auth.your-isolated-instance.example.com/site/oauth2/access_token",
}

c, err := bitbucket.NewOAuthClientCredentialsWithEndpoint(
    "client-id", "client-secret", ep.TokenURL,
)
c, accessToken, err := bitbucket.NewOAuthWithCodeWithEndpoint(
    "client-id", "client-secret", "auth-code", ep,
)
c, accessToken, err := bitbucket.NewOAuthWithRefreshTokenWithEndpoint(
    "client-id", "client-secret", "refresh-token", ep,
)
```

### Custom API base URL (Isolated Cloud Instances, self-hosted, proxies)

When the REST API is reachable under a customer-specific hostname, pass the
base URL alongside the credential:

```go
c, err := bitbucket.NewAPITokenAuthWithBaseUrlStr(
    "you@example.com",
    "your-api-token",
    "https://api.your-isolated-instance.example.com/2.0",
)
```

If the endpoint uses a private CA, supply the PEM bundle:

```go
caBundle, _ := os.ReadFile("/etc/ssl/private-ca.pem")
c, err := bitbucket.NewAPITokenAuthWithBaseUrlStrCaCert(
    "you@example.com",
    "your-api-token",
    "https://api.your-isolated-instance.example.com/2.0",
    caBundle,
)
```

The same `*WithBaseUrlStr` and `*WithBaseUrlStrCaCert` variants exist for
`NewBasicAuth` and `NewOAuthbearerToken`. For OAuth-handshake constructors,
combine the `*WithEndpoint` variant above with `client.SetApiBaseURL` to
redirect REST traffic to the same host. Alternatively, set
`BITBUCKET_API_BASE_URL` in the environment to override the default for any
constructor that does not take a URL argument.

## Usage

### create a pullrequest

```go
package main

import (
        "fmt"

        "github.com/ktrysmt/go-bitbucket"
)

func main() {
        c := bitbucket.NewAPITokenAuth("you@example.com", "your-api-token")

        opt := &bitbucket.PullRequestsOptions{
                Owner:             "your-team",
                RepoSlug:          "awesome-project",
                SourceBranch:      "develop",
                DestinationBranch: "master",
                Title:             "fix bug. #9999",
                CloseSourceBranch: true,
        }

        res, err := c.Repositories.PullRequests.Create(opt)
        if err != nil {
                panic(err)
        }

        fmt.Println(res)
}
```

### create a repository

```go
package main

import (
        "fmt"

        "github.com/ktrysmt/go-bitbucket"
)

func main() {
        c := bitbucket.NewAPITokenAuth("you@example.com", "your-api-token")

        opt := &bitbucket.RepositoryOptions{
                Owner:    "project_name",
                RepoSlug: "repo_name",
                Scm:      "git",
        }

        res, err := c.Repositories.Repository.Create(opt)
        if err != nil {
                panic(err)
        }

        fmt.Println(res)
}
```

## FAQ

### Support Bitbucket API v1.0 ?

It does not correspond yet. Because there are many differences between v2.0 and v1.0.

- Bitbucket API v1.0 <https://confluence.atlassian.com/bitbucket/version-1-423626337.html>

It is officially recommended to use v2.0.
But unfortunately Bitbucket Server (formerly: Stash) API is still v1.0.
And The API v1.0 covers resources that the v2.0 API and API v2.0 is yet to cover.



## License

[Apache License 2.0](./LICENSE)

## Author

[ktrysmt](https://github.com/ktrysmt)
