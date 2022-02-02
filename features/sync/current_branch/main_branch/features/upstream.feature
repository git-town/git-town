Feature: on the main branch with a upstream remote

  Scenario: sync-upstream is set to true
    Given my repo has an upstream repo
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
      |        | git fetch upstream main  |
      |        | git rebase upstream/main |
      |        | git push                 |
      |        | git push --tags          |
      |        | git stash pop            |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, remote, upstream | upstream commit |

  Scenario: sync-upstream is set to false
    Given my repo has an upstream repo
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE         |
      | main   | local    | local commit    |
      |        | remote   | remote commit   |
      |        | upstream | upstream commit |
    And I am on the "main" branch
    And my repo has "git-town.sync-upstream" set to false
    And my workspace has an uncommitted file
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git push --tags          |
      |        | git stash pop            |
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE         |
      | main   | local, remote | remote commit   |
      |        |               | local commit    |
      |        | upstream      | upstream commit |
