@skipWindows
Feature: missing configuration

  Background: running unconfigured
    Given I haven't configured Git Town yet
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
    And I am now on the "feature" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH  | COMMAND               |
      | feature | git checkout main     |
      | main    | git branch -d feature |
    And I am now on the "main" branch
    And Git Town now knows no branch hierarchy
