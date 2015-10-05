Feature: Using Git Town inside a folder that doesn't exist on the main branch

  As a developer using Git Town inside a committed folder that doesn't exist on the main branch
  I want the command to finish properly
  So that my repo is left in a consistent state and I don't lose any data


  Background:
    Given I have feature branches named "current-feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE              | FILE NAME        |
      | main            | local and remote | main commit          | main_file        |
      | current-feature | local and remote | folder commit        | new_folder/file1 |
      | other-feature   | local and remote | other feature commit | file2            |
    And I am on the "current-feature" branch


  Scenario: git-ship
    When I run `git ship -m 'feature done'` in the "new_folder" folder
    Then I am on the "main" branch
    And now I have the following commits
      | BRANCH          | LOCATION         | MESSAGE              | FILE NAME        |
      | main            | local and remote | main commit          | main_file        |
      |                 |                  | feature done         | new_folder/file1 |
      | other-feature   | local and remote | other feature commit | file2            |


  Scenario: git-sync
    When I have an uncommitted file
    And I run `git sync --all` in the "new_folder" folder
    Then it runs the commands
      | BRANCH          | COMMAND                                    |
      | current-feature | git fetch --prune                          |
      | <none>          | cd <%= git_root_folder %>                  |
      | current-feature | git stash -u                               |
      |                 | git checkout main                          |
      | main            | git rebase origin/main                     |
      |                 | git checkout current-feature               |
      | current-feature | git merge --no-edit origin/current-feature |
      |                 | git merge --no-edit main                   |
      |                 | git push                                   |
      |                 | git checkout other-feature                 |
      | other-feature   | git merge --no-edit origin/other-feature   |
      |                 | git merge --no-edit main                   |
      |                 | git push                                   |
      |                 | git checkout current-feature               |
      | current-feature | git push --tags                            |
      |                 | git stash pop                              |
      | <none>          | cd <%= git_folder "new_folder" %>          |
    And I am still on the "current-feature" branch
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                                  | FILE NAME        |
      | main            | local and remote | main commit                              | main_file        |
      | current-feature | local and remote | folder commit                            | new_folder/file1 |
      |                 |                  | main commit                              | main_file        |
      |                 |                  | Merge branch 'main' into current-feature |                  |
      | other-feature   | local and remote | other feature commit                     | file2            |
      |                 |                  | main commit                              | main_file        |
      |                 |                  | Merge branch 'main' into other-feature   |                  |


  Scenario: git-sync --undo
    When I have an uncommitted file
    And I run `git sync --all` in the "new_folder" folder
    And I run `git sync --undo` in the "new_folder" folder
    Then it runs the commands
      | BRANCH          | COMMAND                           |
      | <none>          | cd <%= git_root_folder %>         |
      | current-feature | git stash -u                      |
      |                 | git checkout other-feature        |
      | other-feature   | git checkout current-feature      |
      | current-feature | git checkout main                 |
      | main            | git checkout current-feature      |
      | current-feature | git stash pop                     |
      | <none>          | cd <%= git_folder "new_folder" %> |
    And I am still on the "current-feature" branch
    And I still have my uncommitted file
