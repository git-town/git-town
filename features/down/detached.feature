Feature: move down one position in the detached state

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | beta   | local    | commit 1 |
      |        |          | commit 2 |
    And the current branch is "beta"
    And I ran "git checkout HEAD^"
    When I run "git-town down"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot determine current branch
      """
