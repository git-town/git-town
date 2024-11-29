# bitbucket-app-password

## Bitbucket Cloud

Git Town can interact with Bitbucket Cloud in your name, for example to update
pull requests as branches get created, shipped, or deleted, or to ship pull
requests. To do so, Git Town needs your
[Bitbucket username](bitbucket-username.md) and an
[Bitbucket App Password](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords).

An app password is not the password of your Bitbucket account. It's a special
password that you create so that external applications can interact with
Bitbucket in your name. To create an app password in the Bitbucket web UI, click
on the `Settings` cogwheel, choose `Personal Bitbucket settings`, and then in
the menu on the left `App passwords`. You need to enable these permissions:

- repository: read and write
- pull requests: read and write

## Bitbucket Datacenter

Git Town can interact with Bitbucket Datacenter in your name. To do so, Git Town
needs your [Bitbucket username](bitbucket-username.md) and an
[HTTP access token](https://confluence.atlassian.com/bitbucketserver/http-access-tokens-939515499.html).

An HTTP access token is not the password of your Bitbucket account. It's a
special password that you create so that external applications can interact with
Bitbucket in your name. To create an HTTP access token in the Bitbucket web UI,
click on your Profile picture, choose `Manage account`, and then in the menu on
the left `HTTP access tokens`. You need to enable these permissions:

- Project read
- Repository write

## Setup app password

The best way to enter the Bitbucket app password is via the
[setup assistant](../configuration.md).

## config file

Since your App Password is confidential, you cannot add it to the config file.

## Git metadata

You can configure the App Password manually by running:

```bash
git config [--global] git-town.bitbucket-app-password <token>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
