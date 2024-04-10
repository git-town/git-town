@smoke
Feature: on the main branch

  Background:
    Given the current branch is a feature branch "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    When I run "git-town hack new"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND             |
      | existing | git branch new main |
      |          | git checkout new    |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                     |
      | new      | git checkout main                           |
      | main     | git reset --hard {{ sha 'initial commit' }} |
      |          | git checkout existing                       |
      | existing | git branch -D new                           |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial branches and lineage exist
