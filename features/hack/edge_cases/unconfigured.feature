@skipWindows
Feature: missing configuration

  Background: running unconfigured
    Given Git Town is not configured
    When I run "git-town hack feature" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |

  Scenario: result
    And it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git branch feature main  |
      |        | git checkout feature     |
    And the main branch is now "main"
    And the current branch is now "feature"
    And this branch hierarchy exists now
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout main     |
      | main    | git branch -d feature |
    And the current branch is now "main"
    And no branch hierarchy exists now
