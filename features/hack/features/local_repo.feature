Feature: local repo

  Background:
    Given my repo has a feature branch "existing"
    And my repo does not have an origin remote
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
    And I am on the "existing" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND             |
      | existing | git add -A          |
      |          | git stash           |
      |          | git branch new main |
      |          | git checkout new    |
      | new      | git stash pop       |
    And I am now on the "new" branch
    And my workspace still contains my uncommitted file
    And now these commits exist
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
      | new    | local    | main commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |
