Feature: on the main branch with a upstream remote

  Background:
    Given my repo has an upstream repo

  Scenario: sync-upstream is set to true
    Given my repo contains the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git fetch upstream main  |
      |        | git rebase upstream/main |
      |        | git push                 |
      |        | git push --tags          |
    And all branches are now synchronized
    And I am still on the "main" branch
    And my repo now has the commits
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, remote, upstream | upstream commit |

  Scenario: sync-upstream is set to false
    Given my repo contains the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | local    | local commit    |
      |        | remote   | remote commit   |
      |        | upstream | upstream commit |
    And I am on the "main" branch
    And Git Town's local "git-town.sync-upstream" setting is false
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git push --tags          |
    And all branches are now synchronized
    And I am still on the "main" branch
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE         |
      | main   | local, remote | remote commit   |
      |        |               | local commit    |
      |        | upstream      | upstream commit |
