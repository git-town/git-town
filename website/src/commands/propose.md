# git propose

The _propose_ command helps create a new pull/merge request for the current
feature branch. It opens your code hosting service's website to create a new
proposal in your browser and pre-populates information like branch and
source/target repository. It also [syncs](sync.md) the branch to merge before
opening the pull request.

You can create new pull requests for repositories hosted on:

- [Bitbucket](https://bitbucket.org)
- [Gitea](https://gitea.com)
- [GitHub](https://github.com)
- [GitLab](https://gitlab.com)

### Configuration

You can configure the hosting service type with the
[code-hosting-platform](../preferences/code-hosting-platform.md) setting.

When using SSH identities, this command uses the hostname in the
[code-hosting-origin-hostname](../preferences/code-hosting-origin-hostname.md)
setting.
