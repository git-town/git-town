Feature: observe in detached state

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | branch | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  |
      | branch | local, origin | commit 1 |
      |        | local, origin | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town observe"
  # TODO: fix the broken behavior

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      is local only - branches you want to observe must have a remote branch because they are per definition other people's branches
      """
