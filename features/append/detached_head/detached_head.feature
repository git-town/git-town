Feature: append in detached head state

  Background:
    Given a local Git repo
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      | branch | local    | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town append new"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND             |
      | branch | git checkout -b new |
