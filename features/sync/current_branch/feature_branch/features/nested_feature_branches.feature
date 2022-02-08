Feature: nested feature branches

  Scenario:
    Given my repo has a feature branch "parent"
    And my repo has a feature branch "child" as a child of "parent"
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | remote   | remote main commit   |
      | parent | local    | local parent commit  |
      |        | remote   | remote parent commit |
      | child  | local    | local child commit   |
      |        | remote   | remote child commit  |
    And I am on the "child" branch
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND                           |
      | child  | git fetch --prune --tags          |
      |        | git checkout main                 |
      | main   | git rebase origin/main            |
      |        | git push                          |
      |        | git checkout parent               |
      | parent | git merge --no-edit origin/parent |
      |        | git merge --no-edit main          |
      |        | git push                          |
      |        | git checkout child                |
      | child  | git merge --no-edit origin/child  |
      |        | git merge --no-edit parent        |
      |        | git push                          |
    And all branches are now synchronized
    And I am still on the "child" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local, remote | remote main commit                                       |
      |        |               | local main commit                                        |
      | child  | local, remote | local child commit                                       |
      |        |               | remote child commit                                      |
      |        |               | Merge remote-tracking branch 'origin/child' into child   |
      |        |               | local parent commit                                      |
      |        |               | remote parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | remote main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
      |        |               | Merge branch 'parent' into child                         |
      | parent | local, remote | local parent commit                                      |
      |        |               | remote parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | remote main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
