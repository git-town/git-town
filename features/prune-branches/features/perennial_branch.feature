Feature: remove perennial branch configuration when pruning perennial branches

  Background:
    Given my repo has the perennial branches "active-perennial" and "deleted-perennial"
    And my repo contains the commits
      | BRANCH            | LOCATION      | MESSAGE                  |
      | active-perennial  | local, remote | active-perennial commit  |
      | deleted-perennial | local, remote | deleted-perennial commit |
    And the "deleted-perennial" branch gets deleted on the remote
    And I am on the "deleted-perennial" branch
    And my workspace has an uncommitted file
    When I run "git-town prune-branches"

  Scenario: result
    Then it runs the commands
      | BRANCH            | COMMAND                         |
      | deleted-perennial | git fetch --prune --tags        |
      |                   | git checkout main               |
      | main              | git branch -D deleted-perennial |
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY    | BRANCHES               |
      | local, remote | main, active-perennial |
    And the perennial branches are now "active-perennial"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                           |
      | main   | git branch deleted-perennial {{ sha 'deleted-perennial commit' }} |
      |        | git checkout deleted-perennial                                    |
    And I am now on the "deleted-perennial" branch
    And my workspace still contains my uncommitted file
    And my repo now has the initial branches
    And the perennial branches are now "active-perennial" and "deleted-perennial"
