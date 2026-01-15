@skipWindows
Feature: Create proposals for parked branches

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | parked | local    | parked commit |
    And the current branch is "parked"
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                           |
      | parked | git fetch --prune --tags                                          |
      |        | git merge --no-edit --ff origin/parked                            |
      |        | git push                                                          |
      |        | Finding proposal from parked into main ... ok                     |
      |        | open https://github.com/git-town/git-town/compare/parked?expand=1 |
    And Git Town prints:
      """
      branch parked is no longer parked
      """
    And branch "parked" now has type "feature"
