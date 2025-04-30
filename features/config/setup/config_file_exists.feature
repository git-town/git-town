Feature: don't ask for information already provided by the config file

  @this
  Scenario:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [branches]
      main = "main"
      contribution-regex = "contribute-"
      default-type = "observed"
      feature-regex = "feat-"
      observed-regex = "observed-"
      perennial-regex = "release-"
      perennials = ["staging"]

      [create]
      new-branch-type = "feature"
      share-new-branches = propose

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
      perennial-branches = "fast-forward"
      """
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS  |
      | welcome                     | enter |
      | aliases                     | enter |
      | main branch                 | enter |
      | perennial branches          | enter |
      | perennial regex             | enter |
      | default branch type         | enter |
      | feature regex               | enter |
      | dev-remote                  | enter |
      | forge type                  | enter |
      | origin hostname             | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-prototype-strategy     | enter |
      | sync-upstream               | enter |
      | sync-tags                   | enter |
      | share-new-branches          | enter |
      | push-hook                   | enter |
      | new-branch-type             | enter |
      | ship-strategy               | enter |
      | ship-delete-tracking-branch | enter |
      | save config to config file  | enter |
    Then Git Town runs no commands
    And the main branch is still not set
    And there are still no perennial branches
    And local Git setting "git-town.dev-remote" still doesn't exist
    And local Git setting "git-town.new-branch-type" still doesn't exist
    And local Git setting "git-town.main-branch" still doesn't exist
    And local Git setting "git-town.perennial-branches" still doesn't exist
    And local Git setting "git-town.default-branch-type" still doesn't exist
    And local Git setting "git-town.feature-regex" still doesn't exist
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
