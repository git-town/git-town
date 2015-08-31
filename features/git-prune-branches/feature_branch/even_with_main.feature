Feature: git prune-branches: don't remove the current empty feature branch if there are open changes

  As a developer pruning branches
  I don't want my current empty branch deleted if I have open changes
  So that I can prune my branches without losing current work.


  Background:
    Given the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     | FILE NAME |
      | main   | local and remote | main commit | main_file |
    And I have a stale feature branch named "stale-feature-1" with its tip at "main commit"
    And I have a stale feature branch named "stale-feature-2" with its tip at "main commit"
    And I am on the "stale-feature-1" branch
    And I have an uncommitted file
    When I run `git prune-branches`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                          |
      | stale-feature-1 | git fetch --prune                |
      |                 | git stash -u                     |
      |                 | git checkout main                |
      | main            | git push origin :stale-feature-2 |
      |                 | git branch -d stale-feature-2    |
      |                 | git checkout stale-feature-1     |
      | stale-feature-1 | git stash pop                    |
    And I end up on the "stale-feature-1" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES              |
      | local      | main, stale-feature-1 |
      | remote     | main, stale-feature-1 |
      | coworker   | main                  |


  Scenario: undoing the prune
    When I run `git prune-branches --undo`
    Then it runs the commands
      | BRANCH          | COMMAND                                             |
      | stale-feature-1 | git stash -u                                        |
      |                 | git checkout main                                   |
      | main            | git branch stale-feature-2 <%= sha 'main commit' %> |
      |                 | git push -u origin stale-feature-2                  |
      |                 | git checkout stale-feature-1                        |
      | stale-feature-1 | git stash pop                                       |
    And I end up on the "stale-feature-1" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                               |
      | local      | main, stale-feature-1, stale-feature-2 |
      | remote     | main, stale-feature-1, stale-feature-2 |
      | coworker   | main                                   |
    And I still have my uncommitted file
