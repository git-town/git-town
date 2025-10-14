Feature: ship a parent branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And the current branch is "parent"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                    |
      | parent | git fetch --prune --tags   |
      |        | git checkout main          |
      | main   | git merge --ff-only parent |
      |        | git push                   |
      |        | git push origin :parent    |
      |        | git branch -D parent       |
    And Git Town prints:
      """
      branch "child" is now a child of "main"
      """
    And this lineage exists now
      """
      main
        child
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parent commit |
      | child  | local, origin | child commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git branch parent {{ sha 'parent commit' }} |
      |        | git push -u origin parent                   |
      |        | git checkout parent                         |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And the initial branches and lineage exist now
