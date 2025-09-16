Feature: Cannot create proposals in detached mode

  Background:
    Given a Git repo with origin
    # And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        | local    | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town propose"

  @debug @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | branch | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot propose contribution branches
      """
