#### NAME

repo - view the repository homepage

#### SYNOPSIS

```
git town repo
```

#### DESCRIPTION

Supported only for repositories hosted on [GitHub](http://github.com/),
[GitLab](http://gitlab.com/), and [Bitbucket](https://bitbucket.org/).
When using self-hosted versions this command needs to be configured with
`git config git-town.code-hosting-driver <driver>`
where driver is "github", "gitlab", or "bitbucket".
When using SSH identities, this command needs to be configured with
`git config git-town.code-hosting-origin-hostname <hostname>`
where hostname matches what is in your ssh config file.
