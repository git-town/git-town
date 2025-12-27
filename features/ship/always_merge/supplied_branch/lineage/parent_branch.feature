Feature: ship a parent branch using the always-merge strategy

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
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "child"
    When I run "git-town ship parent" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                            |
      | child  | git fetch --prune --tags           |
      |        | git checkout main                  |
      | main   | git merge --no-ff --edit -- parent |
      |        | git push                           |
      |        | git push origin :parent            |
      |        | git checkout child                 |
      | child  | git branch -D parent               |
    And Git Town prints:
      """
      branch child is now a child of main
      """
    And this lineage exists now
      """
      main
        child
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | main   | local, origin | parent commit         |
      |        |               | Merge branch 'parent' |
      | child  | local, origin | child commit          |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | child  | git branch parent {{ sha 'parent commit' }} |
      |        | git push -u origin parent                   |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE               |
      | main   | local, origin | parent commit         |
      |        |               | Merge branch 'parent' |
      | child  | local, origin | child commit          |
