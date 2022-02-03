Feature: with pull-branch-strategy set to "merge"

  Background:
    Given the pull-branch-strategy configuration is "merge"
    And my repo has a feature branch "feature"
    And my repo contains the commits
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
      | main    | git merge --no-edit origin/main    |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And I am still on the "feature" branch
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, remote | local main commit                                          |
      |         |               | remote main commit                                         |
      |         |               | Merge remote-tracking branch 'origin/main'                 |
      | feature | local, remote | local feature commit                                       |
      |         |               | remote feature commit                                      |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
      |         |               | local main commit                                          |
      |         |               | remote main commit                                         |
      |         |               | Merge remote-tracking branch 'origin/main'                 |
      |         |               | Merge branch 'main' into feature                           |
    And Git Town still has the original branch hierarchy
