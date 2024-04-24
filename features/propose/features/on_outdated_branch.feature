@skipWindows
Feature: sync before proposing

  Background:
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And the current branch is "child"
    And an uncommitted file
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                                                   |
      | child  | git fetch --prune --tags                                                  |
      |        | git add -A                                                                |
      |        | git stash                                                                 |
      |        | git checkout main                                                         |
      | main   | git rebase origin/main                                                    |
      |        | git push                                                                  |
      |        | git checkout parent                                                       |
      | parent | git merge --no-edit origin/parent                                         |
      |        | git merge --no-edit --ff main                                             |
      |        | git push                                                                  |
      |        | git checkout child                                                        |
      | child  | git merge --no-edit origin/child                                          |
      |        | git merge --no-edit parent                                                |
      |        | git push                                                                  |
      |        | git stash pop                                                             |
      | <none> | open https://github.com/git-town/git-town/compare/parent...child?expand=1 |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent...child?expand=1
      """
    And the current branch is still "child"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local, origin | origin main commit                                       |
      |        |               | local main commit                                        |
      | child  | local, origin | local child commit                                       |
      |        |               | origin child commit                                      |
      |        |               | Merge remote-tracking branch 'origin/child' into child   |
      |        |               | local parent commit                                      |
      |        |               | origin parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | origin main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
      |        |               | Merge branch 'parent' into child                         |
      | parent | local, origin | local parent commit                                      |
      |        |               | origin parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | origin main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
