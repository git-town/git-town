Feature: deleting a branch whose parent was shipped and the local branch deleted manually

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
    And origin ships the "parent" branch using the "squash-merge" ship-strategy
    And the current branch is "child"
    And I ran "git branch -d parent"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | child  | git fetch --prune --tags |
      |        | git push origin :child   |
      |        | git checkout main        |
      | main   | git branch -D child      |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | main   | git branch child {{ sha 'child commit' }} |
      |        | git push -u origin child                  |
      |        | git checkout child                        |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
