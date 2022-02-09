Feature: with pull-branch-strategy set to "merge"

  Background:
    Given the "pull-branch-strategy" setting is "merge"
    And my repo has a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git merge --no-edit origin/main    |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And all branches are now synchronized
    And I am still on the "feature" branch
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, origin | local main commit                                          |
      |         |               | origin main commit                                         |
      |         |               | Merge remote-tracking branch 'origin/main'                 |
      | feature | local, origin | local feature commit                                       |
      |         |               | origin feature commit                                      |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
      |         |               | local main commit                                          |
      |         |               | origin main commit                                         |
      |         |               | Merge remote-tracking branch 'origin/main'                 |
      |         |               | Merge branch 'main' into feature                           |
