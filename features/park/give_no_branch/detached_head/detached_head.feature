Feature: cannot park a detached head

  Background:
    Given a local Git repo
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        | local    | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town park"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot park a detached head
      """
    And the initial branches and lineage exist now
