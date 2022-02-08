Feature: sync the current feature branch with a tracking branch

  Background:
    Given my repo has a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | remote   | remote main commit    |
      | feature | local    | local feature commit  |
      |         | remote   | remote feature commit |
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And all branches are now synchronized
    And I am still on the "feature" branch
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, remote | remote main commit                                         |
      |         |               | local main commit                                          |
      | feature | local, remote | local feature commit                                       |
      |         |               | remote feature commit                                      |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
      |         |               | remote main commit                                         |
      |         |               | local main commit                                          |
      |         |               | Merge branch 'main' into feature                           |
