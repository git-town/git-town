# bitbucket-app-password

Git Town can interact with Bitbucket Cloud in your name, for example to update
pull requests as branches get created, shipped, or deleted, or to ship pull
requests. To do so, Git Town needs your
[Bitbucket username](bitbucket-username.md) and an
[Bitbucket App Password](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords).

An app password is not the password of your Bitbucket account. It's a special
password that you create so that external applications can interact with
Bitbucket in your name. To create an app password, click on the `Settings`
cogwheel, choose `Personal Bitbucket settings`, and then in the menu on the left
`App passwords`. You need to enable these permissions:

- repository: read and write
- pull requests: read and write

The best way to enter the Bitbucket app password is via the
[setup assistant](../configuration.md).

## config file

Since your API token is confidential, you cannot add it to the config file.

## Git metadata

You can configure the app password manually by running:

```bash
git config [--global] git-town.bitbucket-app-password <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
