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
      dev-remote = "origin"
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
      | DIALOG                     | KEYS              |
      | welcome                    | enter             |
      | aliases                    | enter             |
      | github connector type: API | enter             |
      | GitHub token               | 1 2 3 4 5 6 enter |
      | token scope: local         | enter             |
      | save config to config file | down enter        |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --local git-town.github-token 123456 |
    And there are still no perennial branches
    And local Git setting "git-town.dev-remote" still doesn't exist
    And local Git setting "git-town.new-branch-type" still doesn't exist
    And local Git setting "git-town.main-branch" still doesn't exist
    And local Git setting "git-town.perennial-branches" still doesn't exist
    And local Git setting "git-town.feature-regex" still doesn't exist
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" is now "123456"
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
    And local Git setting "git-town.unknown-branch-type" still doesn't exist
