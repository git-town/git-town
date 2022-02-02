Feature: restores deleted tracking branch

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
    And the "feature" branch gets deleted on the remote
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git checkout feature       |
      | feature | git merge --no-edit main   |
      |         | git push -u origin feature |
    And I am still on the "feature" branch
    And all branches are now synchronized
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
