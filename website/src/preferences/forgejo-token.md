# Forgejo token

Git Town can interact with Forgejo-based forges (like Codeberg) in your name,
for example to update pull requests as branches get created, shipped, or
deleted. To do so, Git Town needs a personal access token.

To create an API token, follow
[these steps](https://docs.codeberg.org/advanced/access-token) You need an API
token with these permissions:

- repository: read and write

## config file

Since your API token is confidential, you cannot add it to the config file.

## Git metadata

You can configure the API token manually by running:

```wrap
git config [--global] git-town.forgejo-token <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the Forgejo token by setting the `GIT_TOWN_FORGEJO_TOKEN`
environment variable.
