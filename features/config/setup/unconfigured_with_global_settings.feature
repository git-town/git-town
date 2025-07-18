@messyoutput
Feature: setup a new repo when I have configured some things in global Git metadata

  @this
  Scenario:
    Given a Git repo with origin
    And Git Town is not configured
    And global Git setting "git-town.main-branch" is "main"
    And global Git setting "git-town.sync-feature-strategy" is "rebase"
    And global Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And global Git setting "git-town.sync-prototype-strategy" is "compress"
    And global Git setting "git-town.unknown-branch-type" is "observed"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS        |
      | welcome                     | enter       |
      | aliases                     | enter       |
      | main branch                 | enter       |
      | perennial regex             | 1 1 1 enter |
      | feature regex               | 2 2 2 enter |
      | contribution regex          | 3 3 3 enter |
      | observed regex              | 4 4 4 enter |
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
      | new branch type             | enter       |
      | ship strategy               | enter       |
      | ship delete tracking branch | enter       |
      | config storage              | enter       |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.new-branch-type feature          |
      | git config git-town.perennial-regex 111              |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.feature-regex 222                |
      | git config git-town.contribution-regex 333           |
      | git config git-town.observed-regex 444               |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
