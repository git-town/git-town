Feature: already existing remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | origin    |
    When I run "git-town hack existing"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git checkout -b existing |
    And the current branch is now "existing"
    And no commits exist now
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                |
      | existing | git checkout main      |
      | main     | git branch -D existing |
    And the current branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
