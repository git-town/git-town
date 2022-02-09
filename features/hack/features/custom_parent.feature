@skipWindows
Feature: customize the parent for the new feature branch

  Background:
    Given a branch "existing"
    And the current branch is "existing"
    When I run "git-town hack -p new" and answer the prompts:
      | PROMPT                                         | ANSWER        |
      | Please specify the parent branch of 'new'      | [DOWN][ENTER] |
      | Please specify the parent branch of 'existing' | [ENTER]       |

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                     |
      | existing | git fetch --prune --tags    |
      |          | git merge --no-edit main    |
      |          | git push -u origin existing |
      |          | git branch new existing     |
      |          | git checkout new            |
    And the current branch is now "new"
    And Git Town is now aware of this branch hierarchy
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                   |
      | new      | git checkout existing     |
      | existing | git branch -d new         |
      |          | git push origin :existing |
    And the current branch is now "existing"
    And Git Town is now aware of this branch hierarchy
      | BRANCH   | PARENT |
      | existing | main   |
