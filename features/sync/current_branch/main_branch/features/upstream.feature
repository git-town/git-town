Feature: on the main branch with an upstream repo

  Background:
    Given my repo has an upstream repo

  Scenario: sync-upstream is set to true
    Given the commits
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
    And now these commits exist
      | BRANCH | LOCATION                | MESSAGE         |
      | main   | local, origin, upstream | upstream commit |

  Scenario: sync-upstream is set to false
    Given the commits
      | BRANCH | LOCATION | MESSAGE         |
      | main   | local    | local commit    |
      |        | origin   | origin commit   |
      |        | upstream | upstream commit |
    And I am on the "main" branch
    And the "sync-upstream" setting is false
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git rebase origin/main   |
      |        | git push                 |
      |        | git push --tags          |
    And all branches are now synchronized
    And I am still on the "main" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE         |
      | main   | local, origin | origin commit   |
      |        |               | local commit    |
      |        | upstream      | upstream commit |
