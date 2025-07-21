@messyoutput
Feature: setup a new repo when I have configured some things in global Git metadata

  Scenario:
    Given a Git repo with origin
    And Git Town is not configured
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And global Git setting "git-town.feature-regex" is "^feat-"
    And global Git setting "git-town.contribution-regex" is "^cont-"
    And global Git setting "git-town.observed-regex" is "^obs-"
    And global Git setting "git-town.main-branch" is "main"
    And global Git setting "git-town.new-branch-type" is "prototype"
    And global Git setting "git-town.unknown-branch-type" is "prototype"
    And global Git setting "git-town.perennial-branches" is "perennials"
    And global Git setting "git-town.perennial-regex" is "^per-"
    And global Git setting "git-town.push-hook" is "false"
    And global Git setting "git-town.share-new-branches" is "push"
    And global Git setting "git-town.ship-strategy" is "api"
    And global Git setting "git-town.ship-delete-tracking-branch" is "false"
    And global Git setting "git-town.sync-feature-strategy" is "rebase"
    And global Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And global Git setting "git-town.sync-prototype-strategy" is "compress"
    And global Git setting "git-town.sync-tags" is "false"
    And global Git setting "git-town.sync-upstream" is "false"
    And global Git setting "git-town.unknown-branch-type" is "observed"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS        |
      | welcome                     | enter       |
      | aliases                     | enter       |
      | main branch                 | enter       |
      | perennial branches          | space enter |
      | perennial regex             | enter       |
      | feature regex               | enter       |
      | contribution regex          | enter       |
      | observed regex              | enter       |
      | new branch type             | enter       |
      | unknown branch type         | enter       |
      | origin hostname             | enter       |
      | forge type                  | enter       |
      | sync feature strategy       | enter       |
      | sync perennial strategy     | enter       |
      | sync prototype strategy     | enter       |
      | sync upstream               | enter       |
      | sync tags                   | enter       |
      | share new branches          | enter       |
      | push hook                   | enter       |
      | ship strategy               | enter       |
      | ship delete tracking branch | enter       |
      | config storage              | enter       |
    Then Git Town runs the commands
      | COMMAND                                               |
      | git config git-town.main-branch main                  |
      | git config git-town.perennial-branches branch-1       |
      | git config git-town.push-hook false                   |
      | git config git-town.share-new-branches push           |
      | git config git-town.ship-strategy api                 |
      | git config git-town.ship-delete-tracking-branch false |
      | git config git-town.sync-upstream false               |
      | git config git-town.sync-tags false                   |
