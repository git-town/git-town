Feature: already existing remote branch

  Background:
    Given a remote feature branch "existing"
    And an uncommitted file
    When I run "git-town hack existing"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                       |
      | main     | git add -A                    |
      |          | git stash                     |
      |          | git checkout -b existing main |
      | existing | git stash pop                 |
    And the current branch is now "existing"
    And no commits exist now
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
