Feature: nested feature branches

  Scenario:
    Given my repo has a feature branch "parent-feature"
    And my repo has a feature branch "child-feature" as a child of "parent-feature"
    And my repo contains the commits
      | BRANCH         | LOCATION | MESSAGE                      |
      | main           | local    | local main commit            |
      |                | remote   | remote main commit           |
      | parent-feature | local    | local parent feature commit  |
      |                | remote   | remote parent feature commit |
      | child-feature  | local    | local child feature commit   |
      |                | remote   | remote child feature commit  |
    And I am on the "child-feature" branch
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH         | COMMAND                                   |
      | child-feature  | git fetch --prune --tags                  |
      |                | git checkout main                         |
      | main           | git rebase origin/main                    |
      |                | git push                                  |
      |                | git checkout parent-feature               |
      | parent-feature | git merge --no-edit origin/parent-feature |
      |                | git merge --no-edit main                  |
      |                | git push                                  |
      |                | git checkout child-feature                |
      | child-feature  | git merge --no-edit origin/child-feature  |
      |                | git merge --no-edit parent-feature        |
      |                | git push                                  |
    And all branches are now synchronized
    And I am still on the "child-feature" branch
    And my repo now has the commits
      | BRANCH         | LOCATION      | MESSAGE                                                                  |
      | main           | local, remote | remote main commit                                                       |
      |                |               | local main commit                                                        |
      | child-feature  | local, remote | local child feature commit                                               |
      |                |               | remote child feature commit                                              |
      |                |               | Merge remote-tracking branch 'origin/child-feature' into child-feature   |
      |                |               | local parent feature commit                                              |
      |                |               | remote parent feature commit                                             |
      |                |               | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |
      |                |               | remote main commit                                                       |
      |                |               | local main commit                                                        |
      |                |               | Merge branch 'main' into parent-feature                                  |
      |                |               | Merge branch 'parent-feature' into child-feature                         |
      | parent-feature | local, remote | local parent feature commit                                              |
      |                |               | remote parent feature commit                                             |
      |                |               | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |
      |                |               | remote main commit                                                       |
      |                |               | local main commit                                                        |
      |                |               | Merge branch 'main' into parent-feature                                  |
