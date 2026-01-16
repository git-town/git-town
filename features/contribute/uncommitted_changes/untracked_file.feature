Feature: make the current feature branch a contribution branch with an untracked file present

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town contribute"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch feature is now a contribution branch
      """
    And branch "feature" now has type "contribution"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git add -A                  |
      |         | git stash -m "Git Town WIP" |
      |         | git stash pop               |
      |         | git restore --staged .      |
    And branch "feature" now has type "feature"
    And the uncommitted file still exists
