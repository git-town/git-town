@messyoutput
Feature: setup a new repo when I have configured some things in global Git metadata

  @this
  Scenario:
    Given a Git repo with origin
    And Git Town is not configured
    And global Git setting "git-town.feature-regex" is "^kg-"
    And global Git setting "git-town.contribution-regex" is "release-"
    And global Git setting "git-town.observed-regex" is "staging-"
    And global Git setting "git-town.main-branch" is "main"
    And global Git setting "git-town.new-branch-type" is "prototype"
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
      | COMMAND                                            |
      | git config git-town.perennial-regex 111            |
      | git config git-town.unknown-branch-type feature    |
      | git config git-town.feature-regex ^kg-222          |
      | git config git-town.contribution-regex release-333 |
      | git config git-town.observed-regex staging-444     |
