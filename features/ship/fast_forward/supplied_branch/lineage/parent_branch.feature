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
    And the current branch is "child"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship parent"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                    |
      | child  | git fetch --prune --tags   |
      |        | git checkout main          |
      | main   | git merge --ff-only parent |
      |        | git push                   |
      |        | git push origin :parent    |
      |        | git checkout child         |
      | child  | git branch -D parent       |
    And the current branch is now "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | child  | git branch parent {{ sha 'parent commit' }} |
      |        | git push -u origin parent                   |
    And the current branch is now "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parent commit |
      | child  | local, origin | child commit  |
      | parent | local, origin | parent commit |
    And the initial branches and lineage exist
