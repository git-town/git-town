@messyoutput
Feature: setup a new repo when I have configured some things in global Git metadata

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And Git Town is not configured
    And global Git setting "git-town.auto-sync" is "false"
    And global Git setting "git-town.branch-prefix" is "kg-"
    And global Git setting "git-town.contribution-regex" is "^cont-"
    And global Git setting "git-town.detached" is "false"
    And global Git setting "git-town.feature-regex" is "^feat-"
    And global Git setting "git-town.hosting-origin-hostname" is "git"
    And global Git setting "git-town.ignore-uncommitted" is "true"
    And global Git setting "git-town.main-branch" is "main"
    And global Git setting "git-town.new-branch-type" is "prototype"
    And global Git setting "git-town.observed-regex" is "^obs-"
    And global Git setting "git-town.order" is "desc"
    And global Git setting "git-town.perennial-branches" is "perennials"
    And global Git setting "git-town.perennial-regex" is "^per-"
    And global Git setting "git-town.proposal-breadcrumb" is "cli"
    And global Git setting "git-town.proposal-breadcrumb-single" is "false"
    And global Git setting "git-town.push-branches" is "false"
    And global Git setting "git-town.push-hook" is "false"
    And global Git setting "git-town.share-new-branches" is "push"
    And global Git setting "git-town.ship-delete-tracking-branch" is "false"
    And global Git setting "git-town.ship-strategy" is "api"
    And global Git setting "git-town.stash" is "false"
    And global Git setting "git-town.sync-feature-strategy" is "rebase"
    And global Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And global Git setting "git-town.sync-prototype-strategy" is "compress"
    And global Git setting "git-town.sync-tags" is "false"
    And global Git setting "git-town.sync-upstream" is "false"
    And global Git setting "git-town.unknown-branch-type" is "observed"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                      | KEYS        |
      | welcome                     | enter       |
      | aliases                     | enter       |
      | main branch                 | enter       |
      | perennial branches          | space enter |
      | origin hostname             | enter       |
      | forge type                  | enter       |
      | enter all                   | down enter  |
      | perennial regex             | enter       |
      | feature regex               | enter       |
      | contribution regex          | enter       |
      | observed regex              | enter       |
      | branch prefix               | enter       |
      | new branch type             | enter       |
      | unknown branch type         | enter       |
      | sync feature strategy       | enter       |
      | sync perennial strategy     | enter       |
      | sync prototype strategy     | enter       |
      | sync upstream               | enter       |
      | auto sync                   | enter       |
      | sync tags                   | enter       |
      | detached                    | enter       |
      | stash                       | enter       |
      | share new branches          | enter       |
      | push branches               | enter       |
      | push hook                   | enter       |
      | ship strategy               | enter       |
      | ship delete tracking branch | enter       |
      | ignore-uncommitted          | enter       |
      | order                       | enter       |
      | proposal breadcrumb         | enter       |
      | config storage              | enter       |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config git-town.perennial-branches branch-1 |
