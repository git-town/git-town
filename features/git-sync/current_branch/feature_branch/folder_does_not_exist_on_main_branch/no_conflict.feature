Feature: git sync: syncing inside a folder that doesn't exist on the main branch (with open changes)

  As a developer syncing inside a committed folder that doesn't exist on the main branch
  I want the command to at least finish properly
  So that my repo is left in a consistent state and I don't lose any data


  Background:
    Given I have feature branches named "current_feature" and "other_feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE              | FILE NAME        |
      | main            | local and remote | main commit          | main_file        |
      | current_feature | local and remote | folder commit        | new_folder/file1 |
      | other_feature   | local and remote | other feature commit | file2            |
    And I am on the "current_feature" branch
    And I have an uncommitted file
    When I run `git sync --all` in the "new_folder" folder


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                                    |
      | current_feature | git fetch --prune                          |
      |                 | cd <%= git_root_folder %>                  |
      |                 | git stash -u                               |
      |                 | git checkout main                          |
      | main            | git rebase origin/main                     |
      |                 | git checkout current_feature               |
      | current_feature | git merge --no-edit origin/current_feature |
      |                 | git merge --no-edit main                   |
      |                 | git push                                   |
      |                 | git checkout other_feature                 |
      | other_feature   | git merge --no-edit origin/other_feature   |
      |                 | git merge --no-edit main                   |
      |                 | git push                                   |
      |                 | git checkout current_feature               |
      | current_feature | git stash pop                              |
      |                 | cd <%= git_folder "new_folder" %>          |
    And I am still on the "current_feature" branch
    And I still have my uncommitted file
    And now I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                                  | FILE NAME        |
      | main            | local and remote | main commit                              | main_file        |
      | current_feature | local and remote | folder commit                            | new_folder/file1 |
      |                 |                  | main commit                              | main_file        |
      |                 |                  | Merge branch 'main' into current_feature |                  |
      | other_feature   | local and remote | other feature commit                     | file2            |
      |                 |                  | main commit                              | main_file        |
      |                 |                  | Merge branch 'main' into other_feature   |                  |


  Scenario: undo
    When I run `git sync --undo` in the "new_folder" folder
    Then it runs the Git commands
      | BRANCH          | COMMAND                           |
      | current_feature | cd <%= git_root_folder %>         |
      |                 | git stash -u                      |
      |                 | git checkout other_feature        |
      | other_feature   | git checkout current_feature      |
      | current_feature | git checkout main                 |
      | main            | git checkout current_feature      |
      | current_feature | git stash pop                     |
      |                 | cd <%= git_folder "new_folder" %> |
    And I am still on the "current_feature" branch
    And I still have my uncommitted file
