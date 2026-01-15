@skipWindows
Feature: Create proposals for prototype branches

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | prototype | prototype | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE          |
      | main      | local    | main commit      |
      | prototype | local    | prototype commit |
    And the current branch is "prototype"
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                              |
      | prototype | git fetch --prune --tags                                             |
      |           | git push -u origin prototype                                         |
      |           | Finding proposal from prototype into main ... ok                     |
      |           | open https://github.com/git-town/git-town/compare/prototype?expand=1 |
    And Git Town prints:
      """
      branch prototype is no longer a prototype branch
      """
    And branch "prototype" now has type "feature"
