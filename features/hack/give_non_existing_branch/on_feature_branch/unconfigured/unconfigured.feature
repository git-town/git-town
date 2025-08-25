@messyoutput
Feature: missing configuration

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    When I run "git-town hack feature" and enter into the dialog:
      | DIALOG                      | KEYS  |
      | welcome                     | enter |
      | aliases                     | enter |
      | main branch                 | enter |
      | perennial branches          |       |
      | perennial regex             | enter |
      | feature regex               | enter |
      | contribution regex          | enter |
      | observed regex              | enter |
      | new branch type             | enter |
      | unknown branch type         | enter |
      | origin hostname             | enter |
      | forge type                  | enter |
      | sync feature strategy       | enter |
      | sync perennial strategy     | enter |
      | sync prototype strategy     | enter |
      | sync upstream               | enter |
      | sync tags                   | enter |
      | detached                    | enter |
      | stash                       | enter |
      | share new branches          | enter |
      | push hook                   | enter |
      | ship strategy               | enter |
      | ship delete tracking branch | enter |
      | config storage              | enter |

  Scenario: result
    And Git Town runs the commands
      | BRANCH | COMMAND                                              |
      | main   | git fetch --prune --tags                             |
      |        | git config git-town.detached false                   |
      |        | git config git-town.new-branch-type feature          |
      |        | git config git-town.main-branch main                 |
      |        | git config git-town.unknown-branch-type feature      |
      |        | git config git-town.push-hook true                   |
      |        | git config git-town.share-new-branches no            |
      |        | git config git-town.ship-strategy api                |
      |        | git config git-town.ship-delete-tracking-branch true |
      |        | git config git-town.stash true                       |
      |        | git config git-town.sync-feature-strategy merge      |
      |        | git config git-town.sync-perennial-strategy ff-only  |
      |        | git config git-town.sync-prototype-strategy merge    |
      |        | git config git-town.sync-upstream true               |
      |        | git config git-town.sync-tags true                   |
      |        | git checkout -b feature                              |
    And the main branch is now "main"
    And this lineage exists now
      """
      main
        feature
      """
  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout main     |
      | main    | git branch -D feature |
    And no lineage exists now
