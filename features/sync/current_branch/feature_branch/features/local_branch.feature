Feature: syncing the current feature branch without a tracking branch

  Background:
    Given my repo has a local feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      |         | remote   | remote main commit   |
      | feature | local    | local feature commit |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git add -A                 |
      |         | git stash                  |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git push                   |
      |         | git checkout feature       |
      | feature | git merge --no-edit main   |
      |         | git push -u origin feature |
      |         | git stash pop              |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And all branches are now synchronized
