Feature: with upstream remote

  Background:
    Given my repo has an upstream repo
    And my repo has a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | main    | upstream | upstream commit |
      | feature | local    | local commit    |
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git fetch upstream main            |
      |         | git rebase upstream/main           |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And all branches are now synchronized
    And I am still on the "feature" branch
    And my repo now has the commits
      | BRANCH  | LOCATION                | MESSAGE                          |
      | main    | local, remote, upstream | upstream commit                  |
      | feature | local, remote           | local commit                     |
      |         |                         | upstream commit                  |
      |         |                         | Merge branch 'main' into feature |
