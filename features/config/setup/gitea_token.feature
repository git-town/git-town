@messyoutput
Feature: enter the Gitea API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected Gitea platform
    And my repo's "origin" remote is "git@gitea.com:git-town/git-town.git"
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
      | forge type                  | enter             | auto-detect                                 |
      | gitea token                 | 1 2 3 4 5 6 enter |                                             |
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
      | config storage              | enter             | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.gitea-token 123456               |
      | git config git-town.new-branch-type feature          |
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
    And local Git setting "git-town.gitea-token" is now "123456"

  Scenario: select Gitea manually
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                      | DESCRIPTION                                 |
      | welcome                     | enter                     |                                             |
      | aliases                     | enter                     |                                             |
      | main branch                 | enter                     |                                             |
      | perennial branches          |                           | no input here since the dialog doesn't show |
      | perennial regex             | enter                     |                                             |
      | feature regex               | enter                     |                                             |
      | contribution regex          | enter                     |                                             |
      | observed regex              | enter                     |                                             |
      | new branch type             | enter                     |                                             |
      | unknown branch type         | enter                     |                                             |
      | origin hostname             | enter                     |                                             |
      | forge type                  | down down down down enter |                                             |
      | gitea token                 | 1 2 3 4 5 6 enter         |                                             |
      | token scope                 | enter                     |                                             |
      | sync feature strategy       | enter                     |                                             |
      | sync perennial strategy     | enter                     |                                             |
      | sync prototype strategy     | enter                     |                                             |
      | sync upstream               | enter                     |                                             |
      | sync tags                   | enter                     |                                             |
      | share new branches          | enter                     |                                             |
      | push hook                   | enter                     |                                             |
      | ship strategy               | enter                     |                                             |
      | ship delete tracking branch | enter                     |                                             |
      | config storage              | enter                     | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.gitea-token 123456               |
      | git config git-town.new-branch-type feature          |
      | git config git-town.forge-type gitea                 |
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
    And local Git setting "git-town.forge-type" is now "gitea"
    And local Git setting "git-town.gitea-token" is now "123456"

  Scenario: store Gitea API token globally
    And my repo's "origin" remote is "git@gitea.com:git-town/git-town.git"
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
      | gitea token                 | 1 2 3 4 5 6 enter |                                             |
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
      | config storage              | enter             | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.gitea-token 123456      |
      | git config git-town.new-branch-type feature          |
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
    And global Git setting "git-town.gitea-token" is now "123456"

  Scenario: edit global Gitea token
    Given my repo's "origin" remote is "git@gitea.com:git-town/git-town.git"
    And global Git setting "git-town.gitea-token" is "123"
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
      | gitea token                 | backspace backspace backspace 4 5 6 enter |                                             |
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
      | config storage              | enter                                     | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.gitea-token 456         |
      | git config git-town.new-branch-type feature          |
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
    And global Git setting "git-town.gitea-token" is now "456"
