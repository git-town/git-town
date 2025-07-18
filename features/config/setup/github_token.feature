@messyoutput
Feature: enter the GitHub API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected GitHub platform
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | feature regex               | enter             |                                             |
      | contribution regex          | enter             |                                             |
      | observed regex              | enter             |                                             |
      | new branch type             | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type                  | enter             |                                             |
      | github connector type       | enter             |                                             |
      | github token                | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | enter             |                                             |
      | sync feature strategy       | enter             |                                             |
      | sync perennial strategy     | enter             |                                             |
      | sync prototype strategy     | enter             |                                             |
      | sync upstream               | enter             |                                             |
      | sync tags                   | enter             |                                             |
      | share new branches          | enter             |                                             |
      | push hook                   | enter             |                                             |
      | ship strategy               | enter             |                                             |
      | ship delete tracking branch | enter             |                                             |
      | config storage              | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.github-token 123456              |
      | git config git-town.new-branch-type feature          |
      | git config git-town.github-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" is now "123456"

  Scenario: manually selected GitHub
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                           | DESCRIPTION                                 |
      | welcome                     | enter                          |                                             |
      | aliases                     | enter                          |                                             |
      | main branch                 | enter                          |                                             |
      | perennial branches          |                                | no input here since the dialog doesn't show |
      | perennial regex             | enter                          |                                             |
      | feature regex               | enter                          |                                             |
      | contribution regex          | enter                          |                                             |
      | observed regex              | enter                          |                                             |
      | new branch type             | enter                          |                                             |
      | unknown branch type         | enter                          |                                             |
      | origin hostname             | enter                          |                                             |
      | forge type                  | down down down down down enter |                                             |
      | github connector type       | enter                          |                                             |
      | github token                | 1 2 3 4 5 6 enter              |                                             |
      | token scope                 | enter                          |                                             |
      | sync feature strategy       | enter                          |                                             |
      | sync perennial strategy     | enter                          |                                             |
      | sync prototype strategy     | enter                          |                                             |
      | sync upstream               | enter                          |                                             |
      | sync tags                   | enter                          |                                             |
      | share new branches          | enter                          |                                             |
      | push hook                   | enter                          |                                             |
      | ship strategy               | enter                          |                                             |
      | ship delete tracking branch | enter                          |                                             |
      | config storage              | enter                          |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.github-token 123456              |
      | git config git-town.new-branch-type feature          |
      | git config git-town.forge-type github                |
      | git config git-town.github-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.github-token" is now "123456"

  Scenario: remove existing GitHub token
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And local Git setting "git-town.github-token" is "123"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                                | DESCRIPTION                                 |
      | welcome                     | enter                               |                                             |
      | aliases                     | enter                               |                                             |
      | main branch                 | enter                               |                                             |
      | perennial branches          |                                     | no input here since the dialog doesn't show |
      | perennial regex             | enter                               |                                             |
      | feature regex               | enter                               |                                             |
      | contribution regex          | enter                               |                                             |
      | observed regex              | enter                               |                                             |
      | new branch type             | enter                               |                                             |
      | unknown branch type         | enter                               |                                             |
      | origin hostname             | enter                               |                                             |
      | forge type                  | enter                               |                                             |
      | github connector type       | enter                               |                                             |
      | github token                | backspace backspace backspace enter |                                             |
      | sync feature strategy       | enter                               |                                             |
      | sync perennial strategy     | enter                               |                                             |
      | sync prototype strategy     | enter                               |                                             |
      | sync upstream               | enter                               |                                             |
      | sync tags                   | enter                               |                                             |
      | share new branches          | enter                               |                                             |
      | push hook                   | enter                               |                                             |
      | ship strategy               | enter                               |                                             |
      | ship delete tracking branch | enter                               |                                             |
      | config storage              | enter                               |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --unset git-town.github-token             |
      | git config git-town.new-branch-type feature          |
      | git config git-town.github-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist

  Scenario: store GitHub token globally
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | feature regex               | enter             |                                             |
      | contribution regex          | enter             |                                             |
      | observed regex              | enter             |                                             |
      | new branch type             | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type                  | enter             |                                             |
      | github connector type       | enter             |                                             |
      | github token                | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | down enter        |                                             |
      | sync feature strategy       | enter             |                                             |
      | sync perennial strategy     | enter             |                                             |
      | sync prototype strategy     | enter             |                                             |
      | sync upstream               | enter             |                                             |
      | sync tags                   | enter             |                                             |
      | share new branches          | enter             |                                             |
      | push hook                   | enter             |                                             |
      | ship strategy               | enter             |                                             |
      | ship delete tracking branch | enter             |                                             |
      | config storage              | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.github-token 123456     |
      | git config git-town.new-branch-type feature          |
      | git config git-town.github-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And global Git setting "git-town.github-token" is now "123456"

  Scenario: edit global GitHub token
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And global Git setting "git-town.github-token" is "123"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                                      | DESCRIPTION                                 |
      | welcome                     | enter                                     |                                             |
      | aliases                     | enter                                     |                                             |
      | main branch                 | enter                                     |                                             |
      | perennial branches          |                                           | no input here since the dialog doesn't show |
      | perennial regex             | enter                                     |                                             |
      | feature regex               | enter                                     |                                             |
      | contribution regex          | enter                                     |                                             |
      | observed regex              | enter                                     |                                             |
      | new branch type             | enter                                     |                                             |
      | unknown branch type         | enter                                     |                                             |
      | origin hostname             | enter                                     |                                             |
      | forge type                  | enter                                     |                                             |
      | github connector type       | enter                                     |                                             |
      | github token                | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope                 | enter                                     |                                             |
      | sync feature strategy       | enter                                     |                                             |
      | sync perennial strategy     | enter                                     |                                             |
      | sync prototype strategy     | enter                                     |                                             |
      | sync upstream               | enter                                     |                                             |
      | sync tags                   | enter                                     |                                             |
      | share new branches          | enter                                     |                                             |
      | push hook                   | enter                                     |                                             |
      | ship strategy               | enter                                     |                                             |
      | ship delete tracking branch | enter                                     |                                             |
      | config storage              | enter                                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.github-token 456        |
      | git config git-town.new-branch-type feature          |
      | git config git-town.github-connector api             |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And global Git setting "git-town.github-token" is now "456"
