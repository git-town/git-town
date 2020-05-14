Feature: git town-prune-branches: remove perennial branch configuration when pruning perennial branches

  As a developer pruning perennial branches
  I want the configuration to be removed
  So that my config does not contain outdated information.


  Background:
    Given my repository has the branches "active-perennial" and "deleted-perennial"
    And the perennial branches are configured as "active-perennial" and "deleted-perennial"
    And the following commits exist in my repository
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
    And I end up on the "main" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES               |
      | local      | main, active-perennial |
      | remote     | main, active-perennial |
    And the perennial branches are now configured as "active-perennial"


  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                           |
      | main   | git branch deleted-perennial {{ sha 'deleted-perennial commit' }} |
      |        | git checkout deleted-perennial                                    |
    And I end up on the "deleted-perennial" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                                  |
      | local      | main, active-perennial, deleted-perennial |
      | remote     | main, active-perennial                    |
    And the perennial branches are now configured as "active-perennial" and "deleted-perennial"
