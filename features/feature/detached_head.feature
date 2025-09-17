Feature: cannot convert a detached head to a feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        |          | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town feature"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      please check out the branch to make a feature branch
      """
    And the initial branches and lineage exist now
