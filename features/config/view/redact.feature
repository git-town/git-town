Feature: redact API tokens in config output

  Background:
    Given a Git repo with origin
    And Git setting "git-town.bitbucket-app-password" is "bitbucket-password"
    And Git setting "git-town.forgejo-token" is "forgejo-token"
    And Git setting "git-town.gitea-token" is "gitea-token"
    And Git setting "git-town.github-token" is "github-token"
    And Git setting "git-town.gitlab-token" is "gitlab-token"

  Scenario: redaction is the default
    When I run "git-town config"
    Then Git Town prints:
      """
      Hosting:
        browser: (not set), enabled
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

  Scenario: opting out of redaction via --redact=false
    When I run "git-town config --redact=false"
    Then Git Town prints:
      """
      Hosting:
        browser: (not set), enabled
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: bitbucket-password
        Forgejo token: forgejo-token
        Gitea token: gitea-token
        GitHub connector: (not set)
        GitHub token: github-token
        GitLab connector: (not set)
        GitLab token: gitlab-token
      """
