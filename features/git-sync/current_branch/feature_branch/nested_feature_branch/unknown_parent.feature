Feature: git sync: syncing a nested feature branch (without known parent branch)

  As a developer syncing a feature branch without information about its place in the branch hierarchy
  I want to be be able to enter the parent branch efficiently
  So that I am not slowed down much by the process of entering the parent branch.


  Background:
    Given I have a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    |
      | main    | local and remote | main commit    | main_file    |
      | feature | local and remote | feature commit | feature_file |
    And I am on the "feature" branch
    And I have an uncommitted file


  Scenario: choosing the default branch name
    When I run `git sync` and press ENTER
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git stash -u                       |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME    |
      | main    | local and remote | main commit                      | main_file    |
      | feature | local and remote | feature commit                   | feature_file |
      |         |                  | main commit                      | main_file    |
      |         |                  | Merge branch 'main' into feature |              |


  Scenario: entering the number of the parent branch
    When I run `git sync` and enter "1"
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | git stash -u                       |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME    |
      | main    | local and remote | main commit                      | main_file    |
      | feature | local and remote | feature commit                   | feature_file |
      |         |                  | main commit                      | main_file    |
      |         |                  | Merge branch 'main' into feature |              |

