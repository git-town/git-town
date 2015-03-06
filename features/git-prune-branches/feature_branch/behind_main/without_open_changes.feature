Feature: git prune-branches: remove stale feature branches when run on a feature branch (without open changes)

  As a developer pruning branches
  I want all merged branches to be deleted
  So that my remaining branches are relevant and I can focus on my current work.


  Background:
    Given the following commits exist in my repository
      | BRANCH | LOCATION         | MESSAGE     | FILE NAME |
      | main   | local and remote | main commit | main_file |
    And I have a feature branch named "stale_feature_initial" cut from "Initial commit"
    And I have a feature branch named "stale_feature_main" cut from "main commit"
    And I am on the "stale_feature_initial" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH                | COMMAND                                |
      | stale_feature_initial | git fetch --prune                      |
      | stale_feature_initial | git checkout main                      |
      | main                  | git push origin :stale_feature_initial |
      | main                  | git branch -d stale_feature_initial    |
      | main                  | git push origin :stale_feature_main    |
      | main                  | git branch -d stale_feature_main       |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |
      | coworker   | main     |


  Scenario: undoing the operation
    When I run `git prune-branches --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                               |
      | main   | git branch stale_feature_main [SHA:main commit]       |
      | main   | git push -u origin stale_feature_main                 |
      | main   | git branch stale_feature_initial [SHA:Initial commit] |
      | main   | git push -u origin stale_feature_initial              |
      | main   | git checkout stale_feature_initial                    |
    And I end up on the "stale_feature_initial" branch
    Then the existing branches are
        | REPOSITORY | BRANCHES                                        |
        | local      | main, stale_feature_initial, stale_feature_main |
        | remote     | main, stale_feature_initial, stale_feature_main |
        | coworker   | main                                            |
