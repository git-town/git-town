@messyoutput
Feature: remove an existing forge type override

  Background:
    Given a Git repo with origin
    And local Git setting "git-town.forge-type" is "github"
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
      | new branch type             | enter                |                                             |
      | unknown branch type         | enter                |                                             |
      | origin hostname             | enter                |                                             |
      | forge type                  | up up up up up enter |                                             |
      | sync feature strategy       | enter                |                                             |
      | sync perennial strategy     | enter                |                                             |
      | sync prototype strategy     | enter                |                                             |
      | sync upstream               | enter                |                                             |
      | sync tags                   | enter                |                                             |
      | detached                    | enter                |                                             |
      | share new branches          | enter                |                                             |
      | push hook                   | enter                |                                             |
      | ship strategy               | enter                |                                             |
      | ship delete tracking branch | enter                |                                             |
      | config storage              | enter                |                                             |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.detached true                    |
      | git config git-town.new-branch-type feature          |
      | git config --unset git-town.forge-type               |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy ff-only  |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And local Git setting "git-town.forge-type" now doesn't exist

  Scenario: undo
    When I run "git-town undo"
    And local Git setting "git-town.forge-type" is now "github"
