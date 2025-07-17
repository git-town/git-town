@messyoutput
Feature: don't ask for information already provided by the config file

  Scenario:
    Given a Git repo with origin
    And Git Town is not configured
    And the committed configuration file:
      """
      [branches]
      main = "main"
      contribution-regex = "contribute-"
      feature-regex = "feat-"
      observed-regex = "observed-"
      perennial-regex = "release-"
      perennials = ["staging"]
      unknown-type = "observed"

      [create]
      new-branch-type = "feature"
      share-new-branches = "propose"

      [hosting]
      dev-remote = "something"
      origin-hostname = "github.com"
      forge-type = "github"

      [ship]
      delete-tracking-branch = true
      strategy = "api"

      [sync]
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      push-hook = true
      tags = true
      upstream = true

      [sync-strategy]
      feature-branches = "rebase"
      prototype-branches = "merge"
      perennial-branches = "ff-only"
      """
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                | KEYS              |
      | welcome               | enter             |
      | aliases               | enter             |
      | github connector type | enter             |
      | github token          | 1 2 3 4 5 6 enter |
      | token scope           | enter             |
      | config storage        | enter             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config git-town.github-token 123456  |
      | git config git-town.github-connector api |
