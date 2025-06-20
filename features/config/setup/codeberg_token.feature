@messyoutput
Feature: enter the Codeberg API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected Codeberg platform
    And my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | feature regex               | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type: auto-detect     | enter             |                                             |
      | codeberg token              | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | share-new-branches          | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | down enter        |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                           |
      | git config --local git-town.codeberg-token 123456 |
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
      | unknown branch type         | enter                |                                             |
      | dev-remote                  | enter                |                                             |
      | origin hostname             | enter                |                                             |
      | forge type                  | down down down enter |                                             |
      | codeberg token              |    1 2 3 4 5 6 enter |                                             |
      | token scope                 | enter                |                                             |
      | sync-feature-strategy       | enter                |                                             |
      | sync-perennial-strategy     | enter                |                                             |
      | sync-prototype-strategy     | enter                |                                             |
      | sync-upstream               | enter                |                                             |
      | sync-tags                   | enter                |                                             |
      | share-new-branches          | enter                |                                             |
      | push-hook                   | enter                |                                             |
      | new-branch-type             | enter                |                                             |
      | ship-strategy               | enter                |                                             |
      | ship-delete-tracking-branch | enter                |                                             |
      | save config to Git metadata | down enter           |                                             |
    Then Git Town runs the commands
      | COMMAND                                           |
      | git config --local git-town.codeberg-token 123456 |
      | git config git-town.forge-type codeberg           |
    And local Git setting "git-town.forge-type" is now "codeberg"
    And local Git setting "git-town.codeberg-token" is now "123456"

  Scenario: store Codeberge API token globally
    And my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | feature regex               | enter             |                                             |
      | unknown branch type         | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | origin hostname             | enter             |                                             |
      | forge type: auto-detect     | enter             |                                             |
      | codeberg token              | 1 2 3 4 5 6 enter |                                             |
      | token scope                 | down enter        |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | share-new-branches          | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | down enter        |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                            |
      | git config --global git-town.codeberg-token 123456 |
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
      | unknown branch type         | enter                                     |                                             |
      | dev-remote                  | enter                                     |                                             |
      | origin hostname             | enter                                     |                                             |
      | forge type: auto-detect     | enter                                     |                                             |
      | github token                | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope                 | enter                                     |                                             |
      | sync-feature-strategy       | enter                                     |                                             |
      | sync-perennial-strategy     | enter                                     |                                             |
      | sync-prototype-strategy     | enter                                     |                                             |
      | sync-upstream               | enter                                     |                                             |
      | sync-tags                   | enter                                     |                                             |
      | share-new-branches          | enter                                     |                                             |
      | push-hook                   | enter                                     |                                             |
      | new-branch-type             | down enter                                |                                             |
      | ship-strategy               | enter                                     |                                             |
      | ship-delete-tracking-branch | enter                                     |                                             |
      | save config to Git metadata | down enter                                |                                             |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --global git-town.codeberg-token 456 |
    And global Git setting "git-town.codeberg-token" is now "456"
