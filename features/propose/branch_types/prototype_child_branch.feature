@skipWindows
Feature: Create proposals for prototype branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | parent    | feature   | main   | local     |
      | prototype | prototype | parent | local     |
    And the current branch is "prototype"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exist
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                       |
      | prototype | git fetch --prune --tags                                                      |
      |           | Looking for proposal online ... ok                                            |
      | prototype | git checkout parent                                                           |
      | parent    | git push -u origin parent                                                     |
      |           | git checkout prototype                                                        |
      | prototype | git push -u origin prototype                                                  |
      |           | open https://github.com/git-town/git-town/compare/parent...prototype?expand=1 |
    And Git Town prints:
      """
      branch "prototype" is no longer a prototype branch
      """
    And branch "prototype" now has type "feature"
