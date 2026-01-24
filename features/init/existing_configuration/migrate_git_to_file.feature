@messyoutput
Feature: migrate existing configuration in Git metadata to a config file

  Background:
    Given a Git repo with origin
    And the main branch is "main"
    And local Git setting "git-town.auto-sync" is "false"
    And local Git setting "git-town.branch-prefix" is "acme-"
    And local Git setting "git-town.contribution-regex" is "coworker-.*"
    And local Git setting "git-town.dev-remote" is "fork"
    And local Git setting "git-town.feature-regex" is "user-.*"
    And local Git setting "git-town.forge-type" is "github"
    And local Git setting "git-town.github-connector" is "api"
    And local Git setting "git-town.hosting-origin-hostname" is "github.example.com"
    And local Git setting "git-town.ignore-uncommitted" is "true"
    And local Git setting "git-town.new-branch-type" is "prototype"
    And local Git setting "git-town.observed-regex" is "other-.*"
    And local Git setting "git-town.order" is "desc"
    And local Git setting "git-town.perennial-branches" is "qa"
    And local Git setting "git-town.perennial-regex" is "release-.*"
    And local Git setting "git-town.proposal-breadcrumb" is "cli"
    And local Git setting "git-town.push-branches" is "true"
    And local Git setting "git-town.push-hook" is "true"
    And local Git setting "git-town.share-new-branches" is "no"
    And local Git setting "git-town.ship-delete-tracking-branch" is "false"
    And local Git setting "git-town.ship-strategy" is "squash-merge"
    And local Git setting "git-town.stash" is "false"
    And local Git setting "git-town.sync-feature-strategy" is "merge"
    And local Git setting "git-town.sync-perennial-strategy" is "rebase"
    And local Git setting "git-town.sync-tags" is "false"
    And local Git setting "git-town.sync-upstream" is "true"
    And local Git setting "git-town.unknown-branch-type" is "observed"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                      | KEYS       |
      | welcome                     | enter      |
      | aliases                     | enter      |
      | main branch                 | enter      |
      | perennial branches          | enter      |
      | dev-remote                  | enter      |
      | origin hostname             | enter      |
      | forge type                  | enter      |
      | github connector            | enter      |
      | github token                | enter      |
      | enter all                   | down enter |
      | perennial regex             | enter      |
      | feature regex               | enter      |
      | contribution regex          | enter      |
      | observed regex              | enter      |
      | branch prefix               | enter      |
      | new branch type             | enter      |
      | unknown branch type         | enter      |
      | sync feature strategy       | enter      |
      | sync perennial strategy     | enter      |
      | sync prototype strategy     | enter      |
      | sync upstream               | enter      |
      | auto sync                   | enter      |
      | sync tags                   | enter      |
      | detached                    | enter      |
      | stash                       | enter      |
      | share new branches          | enter      |
      | push branches               | enter      |
      | push hook                   | enter      |
      | ship strategy               | enter      |
      | ship delete tracking branch | enter      |
      | ignore-uncommitted          | enter      |
      | order                       | enter      |
      | proposal breadcrumb         | enter      |
      | proposal breadcrumb single  | enter      |
      | config storage              | down enter |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                                 |
      | git config --unset git-town.auto-sync                   |
      | git config --unset git-town.branch-prefix               |
      | git config --unset git-town.contribution-regex          |
      | git config --unset git-town.dev-remote                  |
      | git config --unset git-town.feature-regex               |
      | git config --unset git-town.forge-type                  |
      | git config --unset git-town.github-connector            |
      | git config --unset git-town.hosting-origin-hostname     |
      | git config --unset git-town.ignore-uncommitted          |
      | git config --unset git-town.main-branch                 |
      | git config --unset git-town.new-branch-type             |
      | git config --unset git-town.observed-regex              |
      | git config --unset git-town.order                       |
      | git config --unset git-town.perennial-regex             |
      | git config --unset git-town.proposal-breadcrumb         |
      | git config --unset git-town.push-branches               |
      | git config --unset git-town.push-hook                   |
      | git config --unset git-town.share-new-branches          |
      | git config --unset git-town.ship-delete-tracking-branch |
      | git config --unset git-town.ship-strategy               |
      | git config --unset git-town.stash                       |
      | git config --unset git-town.sync-feature-strategy       |
      | git config --unset git-town.sync-perennial-strategy     |
      | git config --unset git-town.sync-tags                   |
      | git config --unset git-town.sync-upstream               |
      | git config --unset git-town.unknown-branch-type         |
      | git config --unset git-town.perennial-branches          |
    And local Git setting "git-town.auto-sync" now doesn't exist
    And local Git setting "git-town.contribution-regex" now doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-connector" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.ignore-uncommitted" now doesn't exist
    And local Git setting "git-town.new-branch-type" now doesn't exist
    And local Git setting "git-town.observed-regex" now doesn't exist
    And local Git setting "git-town.order" now doesn't exist
    And local Git setting "git-town.perennial-regex" now doesn't exist
    And local Git setting "git-town.proposal-breadcrumb" now doesn't exist
    And local Git setting "git-town.push-branches" now doesn't exist
    And local Git setting "git-town.push-hook" now doesn't exist
    And local Git setting "git-town.share-new-branches" now doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" now doesn't exist
    And local Git setting "git-town.ship-strategy" now doesn't exist
    And local Git setting "git-town.stash" now doesn't exist
    And local Git setting "git-town.sync-feature-strategy" now doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" now doesn't exist
    And local Git setting "git-town.sync-tags" now doesn't exist
    And local Git setting "git-town.sync-upstream" now doesn't exist
    And local Git setting "git-town.unknown-branch-type" now doesn't exist
    And the configuration file is now:
      """
      # See https://www.git-town.com/configuration-file for details

      [branches]
      contribution-regex = "coworker-.*"
      feature-regex = "user-.*"
      main = "main"
      observed-regex = "other-.*"
      order = "desc"
      perennials = ["qa"]
      perennial-regex = "release-.*"
      unknown-type = "observed"

      [create]
      branch-prefix = "acme-"
      new-branch-type = "prototype"
      share-new-branches = "no"
      stash = false

      [hosting]
      dev-remote = "fork"
      forge-type = "github"
      github-connector = "api"
      origin-hostname = "github.example.com"

      [propose]
      breadcrumb = "cli"
      breadcrumb-single = true

      [ship]
      delete-tracking-branch = false
      ignore-uncommitted = true
      strategy = "squash-merge"

      [sync]
      auto-sync = false
      detached = false
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      prototype-strategy = "merge"
      push-branches = true
      push-hook = true
      tags = false
      upstream = true
      """
    And the main branch is now not set
    And there are now no perennial branches

  Scenario: undo
    When I run "git-town undo"
    Then local Git setting "git-town.auto-sync" is now "false"
    And local Git setting "git-town.contribution-regex" is now "coworker-.*"
    And local Git setting "git-town.dev-remote" is now "fork"
    And local Git setting "git-town.feature-regex" is now "user-.*"
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.github-connector" is now "api"
    And local Git setting "git-town.hosting-origin-hostname" is now "github.example.com"
    And local Git setting "git-town.ignore-uncommitted" is now "true"
    And local Git setting "git-town.new-branch-type" is now "prototype"
    And local Git setting "git-town.observed-regex" is now "other-.*"
    And local Git setting "git-town.order" is now "desc"
    And local Git setting "git-town.perennial-regex" is now "release-.*"
    And local Git setting "git-town.proposal-breadcrumb" is now "cli"
    And local Git setting "git-town.push-branches" is now "true"
    And local Git setting "git-town.push-hook" is now "true"
    And local Git setting "git-town.share-new-branches" is now "no"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "false"
    And local Git setting "git-town.ship-strategy" is now "squash-merge"
    And local Git setting "git-town.stash" is now "false"
    And local Git setting "git-town.sync-feature-strategy" is now "merge"
    And local Git setting "git-town.sync-perennial-strategy" is now "rebase"
    And local Git setting "git-town.sync-tags" is now "false"
    And local Git setting "git-town.sync-upstream" is now "true"
    And local Git setting "git-town.unknown-branch-type" is now "observed"
    And the main branch is now "main"
