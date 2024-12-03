Feature: deleting a branch that conflicts with the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME        | FILE CONTENT   |
      | main   | local, origin | main commit   | conflicting_file | main content   |
      | parent | local, origin | parent commit | parent_file      | parent content |
      | child  | local, origin | child commit  | conflicting_file | child content  |
    And the current branch is "child"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | child  | git fetch --prune --tags |
      |        | git push origin :child   |
      |        | git checkout main        |
      | main   | git branch -D child      |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, parent |
    And this lineage exists now
      | BRANCH | PARENT |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | main   | git branch child {{ sha 'child commit' }} |
      |        | git push -u origin child                  |
      |        | git checkout child                        |
    And the current branch is still "child"
    And the branches are now
      | REPOSITORY    | BRANCHES            |
      | local, origin | main, child, parent |
    And the initial lineage exists now
