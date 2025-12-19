@messyoutput
Feature: don't ask for information already provided by the config file

  Scenario:
    Given a Git repo with origin
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
      stash = true

      [hosting]
      dev-remote = "something"
      origin-hostname = "github.com"
      forge-type = "github"
      github-connector = "gh"

      [ship]
      delete-tracking-branch = true
      strategy = "api"

      [sync]
      detached = false
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
    And Git Town is not configured
    And global Git setting "git-town.github-token" is "123456"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG             | KEYS  |
      | welcome            | enter |
      | aliases            | enter |
      | perennial branches | enter |
      | enter all          | enter |
      | config storage     | enter |
    Then Git Town runs no commands
    And global Git setting "git-town.github-token" is still "123456"
    # keep-sorted start
    And local Git setting "git-town.contribution-regex" still doesn't exist
    And local Git setting "git-town.dev-remote" still doesn't exist
    And local Git setting "git-town.feature-regex" still doesn't exist
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" still doesn't exist
    And local Git setting "git-town.main-branch" still doesn't exist
    And local Git setting "git-town.new-branch-type" still doesn't exist
    And local Git setting "git-town.perennial-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.stash" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.unknown-branch-type" still doesn't exist
    # keep-sorted end
    And there are still no perennial branches
