Feature: already existing unfetched remote branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | origin    |
    And an uncommitted file
    When I run "git-town append existing"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | main     | git add -A               |
      |          | git stash                |
      |          | git checkout -b existing |
      | existing | git stash pop            |
    And the current branch is now "existing"
    And the initial commits exist
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                |
      | existing | git add -A             |
      |          | git stash              |
      |          | git checkout main      |
      | main     | git branch -D existing |
      |          | git stash pop          |
    And the current branch is now "main"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
