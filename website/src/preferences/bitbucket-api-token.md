# bitbucket-api-token

## Bitbucket Cloud

Git Town can interact with Bitbucket Cloud in your name,
for example to update pull requests as branches get created, shipped,
or deleted, or to ship pull requests.
To do so, Git Town needs your [Bitbucket username](bitbucket-username.md)
(the email address of your Atlassian account) and an
[API token with scopes](https://support.atlassian.com/bitbucket-cloud/docs/create-an-api-token).

An API token is not the password of your Atlassian account.
It's a special credential that you create so
that external applications can interact with Bitbucket in your name.
To create an API token, open the
[security settings of your Atlassian account](https://id.atlassian.com/manage-profile/security/api-tokens),
click `Create API token with scopes`, choose `Bitbucket` as the app,
and select these scopes:

- `read:user:bitbucket`
- `read:repository:bitbucket`
- `read:pullrequest:bitbucket`
- `write:pullrequest:bitbucket`

If you still have a Bitbucket app password configured,
Git Town automatically renames the deprecated `git-town.bitbucket-app-password`
setting to `git-town.bitbucket-api-token`.
Bitbucket no longer accepts app passwords,
so please replace the old value with an API token.

## Bitbucket Data Center

Git Town can interact with Bitbucket Data Center in your name.
To do so, Git Town needs your [Bitbucket username](bitbucket-username.md) and an
[HTTP access token](https://confluence.atlassian.com/bitbucketserver/http-access-tokens-939515499.html).

An HTTP access token is not the password of your Bitbucket account.
It's a special password that you create so
that external applications can interact with Bitbucket in your name.
To create an HTTP access token in the Bitbucket web UI,
click on your Profile picture, choose `Manage account`,
and then in the menu on the left `HTTP access tokens`.
You need to enable these permissions:

- Project read
- Repository write

## config file

Since your API token or access token is confidential,
you cannot add it to the config file.

## Git metadata

You can configure the API token or access token manually by running:

```wrap
git config [--global] git-town.bitbucket-api-token <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine.
Without it, the setting applies only to the current repository.

## environment variable

You can configure the Bitbucket API token by setting the
`GIT_TOWN_BITBUCKET_API_TOKEN` environment variable.
