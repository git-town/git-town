Feature: redact API tokens in config output

  Background:
    Given a Git repo with origin
    And Git setting "git-town.bitbucket-app-password" is "bitbucket-password"
    And Git setting "git-town.forgejo-token" is "forgejo-token"
    And Git setting "git-town.gitea-token" is "gitea-token"
    And Git setting "git-town.github-token" is "github-token"
    And Git setting "git-town.gitlab-token" is "gitlab-token"
    When I run "git-town config --redact"

  Scenario: result
    Then Git Town prints:
      """
      Hosting:
        browser: (not set)
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (configured)
        Forgejo token: (configured)
        Gitea token: (configured)
        GitHub connector: (not set)
        GitHub token: (configured)
        GitLab connector: (not set)
        GitLab token: (configured)
      """
