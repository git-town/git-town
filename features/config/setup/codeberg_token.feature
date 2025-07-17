@messyoutput
Feature: enter the Codeberg API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected Codeberg platform
    Given my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
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
      | unknown branch type         | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type                  | enter             |                                             |
      | codeberg token              | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | enter             |                                             |
      | sync feature strategy       | enter             |                                             |
      | sync perennial strategy     | enter             |                                             |
      | sync prototype strategy     | enter             |                                             |
      | sync upstream               | enter             |                                             |
      | sync tags                   | enter             |                                             |
      | share new branches          | enter             |                                             |
      | push hook                   | enter             |                                             |
      | new branch type             | enter             |                                             |
      | ship strategy               | enter             |                                             |
      | ship delete tracking branch | enter             |                                             |
      | config storage              | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.codeberg-token 123456            |
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
    And local Git setting "git-town.codeberg-token" is now "123456"

  Scenario: select Codeberg manually
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                 | DESCRIPTION                                 |
      | welcome                     | enter                |                                             |
      | aliases                     | enter                |                                             |
      | main branch                 | enter                |                                             |
      | perennial branches          |                      | no input here since the dialog doesn't show |
      | perennial regex             | enter                |                                             |
      | feature regex               | enter                |                                             |
      | contribution regex          | enter                |                                             |
      | observed regex              | enter                |                                             |
      | unknown branch type         | enter                |                                             |
      | origin hostname             | enter                |                                             |
      | forge type                  | down down down enter |                                             |
      | codeberg token              | 1 2 3 4 5 6 enter    |                                             |
      | token scope                 | enter                |                                             |
      | sync feature strategy       | enter                |                                             |
      | sync perennial strategy     | enter                |                                             |
      | sync prototype strategy     | enter                |                                             |
      | sync upstream               | enter                |                                             |
      | sync tags                   | enter                |                                             |
      | share new branches          | enter                |                                             |
      | push hook                   | enter                |                                             |
      | new branch type             | enter                |                                             |
      | ship strategy               | enter                |                                             |
      | ship delete tracking branch | enter                |                                             |
      | config storage              | enter                |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.codeberg-token 123456            |
      | git config git-town.new-branch-type feature          |
      | git config git-town.forge-type codeberg              |
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
    And local Git setting "git-town.forge-type" is now "codeberg"
    And local Git setting "git-town.codeberg-token" is now "123456"

  Scenario: store Codeberge API token globally
    And my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main-branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | feature regex               | enter             |                                             |
      | contribution regex          | enter             |                                             |
      | observed regex              | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type                  | enter             |                                             |
      | codeberg token              | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | down enter        |                                             |
      | sync feature strategy       | enter             |                                             |
      | sync perennial strategy     | enter             |                                             |
      | sync prototype strategy     | enter             |                                             |
      | sync upstream               | enter             |                                             |
      | sync tags                   | enter             |                                             |
      | share new branches          | enter             |                                             |
      | push hook                   | enter             |                                             |
      | new branch type             | enter             |                                             |
      | ship strategy               | enter             |                                             |
      | ship delete tracking branch | enter             |                                             |
      | config storage              | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.codeberg-token 123456   |
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
    And global Git setting "git-town.codeberg-token" is now "123456"

  Scenario: edit global Codeberge API token
    And my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
    Given global Git setting "git-town.codeberg-token" is "123"
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
      | unknown branch type         | enter                                     |                                             |
      | origin hostname             | enter                                     |                                             |
      | forge type                  | enter                                     |                                             |
      | codeberg token              | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope                 | enter                                     |                                             |
      | sync feature strategy       | enter                                     |                                             |
      | sync perennial strategy     | enter                                     |                                             |
      | sync prototype strategy     | enter                                     |                                             |
      | sync upstream               | enter                                     |                                             |
      | sync tags                   | enter                                     |                                             |
      | share new branches          | enter                                     |                                             |
      | push hook                   | enter                                     |                                             |
      | new branch type             | enter                                     |                                             |
      | ship strategy               | enter                                     |                                             |
      | ship delete tracking branch | enter                                     |                                             |
      | config storage              | enter                                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global git-town.codeberg-token 456      |
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
    And global Git setting "git-town.codeberg-token" is now "456"
