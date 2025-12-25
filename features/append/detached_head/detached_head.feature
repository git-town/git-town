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

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      please check out the branch to which you want to append a child
      """
