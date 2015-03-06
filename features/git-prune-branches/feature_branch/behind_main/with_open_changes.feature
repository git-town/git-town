Feature: git prune-branches: don't remove the current empty feature branch if there are open changes

  As a developer pruning branches
  I don't want my current empty branch deleted if I have open changes
  So that I can prune my branches without losing current work.


  Background:
    Given the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     | FILE NAME |
      | main   | local and remote | main commit | main_file |
    And I have a feature branch cut from "Initial commit" named "stale_feature_initial"
    And I have a feature branch cut from "main commit" named "stale_feature_main"
    And I am on the "stale_feature_initial" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH                | COMMAND                             |
      | stale_feature_initial | git fetch --prune                   |
      | stale_feature_initial | git stash -u                        |
      | stale_feature_initial | git checkout main                   |
      | main                  | git push origin :stale_feature_main |
      | main                  | git branch -d stale_feature_main    |
      | main                  | git checkout stale_feature_initial  |
      | stale_feature_initial | git stash pop                       |
    And I end up on the "stale_feature_initial" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                    |
      | local      | main, stale_feature_initial |
      | remote     | main, stale_feature_initial |
      | coworker   | main                        |


  Scenario: undoing the prune
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH                | COMMAND                                            |
      | stale_feature_initial | git stash -u                                       |
      | stale_feature_initial | git checkout main                                  |
      | main                  | git branch stale_feature_main [SHA:main commit] |
      | main                  | git push -u origin stale_feature_main              |
      | main                  | git checkout stale_feature_initial                 |
      | stale_feature_initial | git stash pop                                      |
    And I end up on the "stale_feature_initial" branch
    Then the existing branches are
      | REPOSITORY | BRANCHES                                        |
      | local      | main, stale_feature_main, stale_feature_initial |
      | remote     | main, stale_feature_main, stale_feature_initial |
      | coworker   | main                                            |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
