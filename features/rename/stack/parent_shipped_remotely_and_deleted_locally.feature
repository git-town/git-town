Feature: renaming a branch whose parent was shipped and the local branch deleted manually

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
    When I run "git-town rename new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | child  | git fetch --prune --tags    |
      |        | git branch --move child new |
      |        | git checkout new            |
      | new    | git push -u origin new      |
      |        | git push origin :child      |
    And this lineage exists now
      """
      main
        new
      """
    And the branches are now
      | REPOSITORY    | BRANCHES  |
      | local, origin | main, new |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | new    | git branch child {{ sha 'child commit' }} |
      |        | git push -u origin child                  |
      |        | git checkout child                        |
      | child  | git branch -D new                         |
      |        | git push origin :new                      |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, child |
