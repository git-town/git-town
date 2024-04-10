Feature: already existing remote branch

  Background:
    Given a remote feature branch "existing"
    And an uncommitted file
    When I run "git-town hack existing"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | main     | git add -A               |
      |          | git stash                |
      |          | git branch existing main |
      |          | git checkout existing    |
      | existing | git stash pop            |
    And the current branch is now "existing"
    And no commits exist now
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
    And the uncommitted file still exists
