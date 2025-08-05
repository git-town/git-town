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
      | dev remote                  | enter |
      | origin hostname             | enter |
      | forge type                  | enter |
      | sync feature strategy       | enter |
      | sync perennial strategy     | enter |
      | sync prototype strategy     | enter |
      | sync upstream               | enter |
      | sync tags                   | enter |
      | share new branches          | enter |
      | push hook                   | enter |
      | ship strategy               | enter |
      | ship delete tracking branch | enter |
      | config storage              | enter |

  Scenario: result
    And Git Town runs the commands
      | BRANCH | COMMAND                              |
      | main   | git fetch --prune --tags             |
      |        | git config git-town.main-branch main |
      |        | git checkout -b feature              |
    And the main branch is now "main"
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout main     |
      | main    | git branch -D feature |
    And no lineage exists now
