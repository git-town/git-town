Feature: git sync: syncing a nested feature branch (without known parent branch)

  As a developer syncing a feature branch without information about its place in the branch hierarchy
  I want to be be able to enter the parent branch efficiently
  So that I am not slowed down much by the process of entering the parent branch.


  Background:
    Given I have feature branches named "feature-1" and "feature-2"
    And Git Town has no branch hierarchy information for "feature-1" and "feature-2"
    And the following commits exist in my repository
      | BRANCH    | LOCATION         | MESSAGE          | FILE NAME      |
      | main      | local and remote | main commit      | main_file      |
      | feature-1 | local and remote | feature 1 commit | feature_1_file |
      | feature-2 | local and remote | feature 2 commit | feature_2_file |
    And I am on the "feature-2" branch
    And I have an uncommitted file


  Scenario: choosing the default branch name
    When I run `git sync` and enter ""
    Then it runs the Git commands
      | BRANCH    | COMMAND                              |
      | feature-2 | git fetch --prune                    |
      |           | git stash -u                         |
      |           | git checkout main                    |
      | main      | git rebase origin/main               |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit origin/feature-2 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git stash pop                        |
    And I am still on the "feature-2" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION         | MESSAGE                            | FILE NAME      |
      | main      | local and remote | main commit                        | main_file      |
      | feature-1 | local and remote | feature 1 commit                   | feature_1_file |
      | feature-2 | local and remote | feature 2 commit                   | feature_2_file |
      |           |                  | main commit                        | main_file      |
      |           |                  | Merge branch 'main' into feature-2 |                |


  Scenario: entering the number of the master branch
    When I run `git sync` and enter "1"
    Then it runs the Git commands
      | BRANCH    | COMMAND                              |
      | feature-2 | git fetch --prune                    |
      |           | git stash -u                         |
      |           | git checkout main                    |
      | main      | git rebase origin/main               |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit origin/feature-2 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git stash pop                        |
    And I am still on the "feature-2" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION         | MESSAGE                            | FILE NAME      |
      | main      | local and remote | main commit                        | main_file      |
      | feature-1 | local and remote | feature 1 commit                   | feature_1_file |
      | feature-2 | local and remote | feature 2 commit                   | feature_2_file |
      |           |                  | main commit                        | main_file      |
      |           |                  | Merge branch 'main' into feature-2 |                |


  Scenario: entering the number of another branch
    When I run `git sync` and enter "2"
    Then it runs the Git commands
      | BRANCH    | COMMAND                              |
      | feature-2 | git fetch --prune                    |
      |           | git stash -u                         |
      |           | git checkout main                    |
      | main      | git rebase origin/main               |
      |           | git checkout feature-1               |
      | feature-1 | git merge --no-edit origin/feature-1 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit origin/feature-2 |
      |           | git merge --no-edit feature-1        |
      |           | git push                             |
      |           | git stash pop                        |
    And I am still on the "feature-2" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION         | MESSAGE                                 | FILE NAME      |
      | main      | local and remote | main commit                             | main_file      |
      | feature-1 | local and remote | feature 1 commit                        | feature_1_file |
      |           |                  | main commit                             | main_file      |
      |           |                  | Merge branch 'main' into feature-1      |                |
      | feature-2 | local and remote | feature 2 commit                        | feature_2_file |
      |           |                  | feature 1 commit                        | feature_1_file |
      |           |                  | main commit                             | main_file      |
      |           |                  | Merge branch 'main' into feature-1      |                |
      |           |                  | Merge branch 'feature-1' into feature-2 |                |


  Scenario: entering the name of the master branch
    When I run `git sync` and enter "main"
    Then it runs the Git commands
      | BRANCH    | COMMAND                              |
      | feature-2 | git fetch --prune                    |
      |           | git stash -u                         |
      |           | git checkout main                    |
      | main      | git rebase origin/main               |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit origin/feature-2 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git stash pop                        |
    And I am still on the "feature-2" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION         | MESSAGE                            | FILE NAME      |
      | main      | local and remote | main commit                        | main_file      |
      | feature-1 | local and remote | feature 1 commit                   | feature_1_file |
      | feature-2 | local and remote | feature 2 commit                   | feature_2_file |
      |           |                  | main commit                        | main_file      |
      |           |                  | Merge branch 'main' into feature-2 |                |


  Scenario: entering the name of another branch
    When I run `git sync` and enter "feature-1"
    Then it runs the Git commands
      | BRANCH    | COMMAND                              |
      | feature-2 | git fetch --prune                    |
      |           | git stash -u                         |
      |           | git checkout main                    |
      | main      | git rebase origin/main               |
      |           | git checkout feature-1               |
      | feature-1 | git merge --no-edit origin/feature-1 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit origin/feature-2 |
      |           | git merge --no-edit feature-1        |
      |           | git push                             |
      |           | git stash pop                        |
    And I am still on the "feature-2" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH    | LOCATION         | MESSAGE                                 | FILE NAME      |
      | main      | local and remote | main commit                             | main_file      |
      | feature-1 | local and remote | feature 1 commit                        | feature_1_file |
      |           |                  | main commit                             | main_file      |
      |           |                  | Merge branch 'main' into feature-1      |                |
      | feature-2 | local and remote | feature 2 commit                        | feature_2_file |
      |           |                  | feature 1 commit                        | feature_1_file |
      |           |                  | main commit                             | main_file      |
      |           |                  | Merge branch 'main' into feature-1      |                |
      |           |                  | Merge branch 'feature-1' into feature-2 |                |
