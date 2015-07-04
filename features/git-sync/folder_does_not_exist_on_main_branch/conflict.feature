Feature: git sync: syncing inside a folder that doesn't exist on the main branch

  (see ./no_conflict.feature)


  Background:
    Given I have feature branches named "current_feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main            | local and remote | conflicting main commit    | conflicting_file | main content    |
      | current_feature | local            | conflicting feature commit | conflicting_file | feature content |
      |                 |                  | folder commit              | new_folder/file1 |                 |
      | other-feature   | local and remote | other feature commit       | file2            |                 |
    And I am on the "current_feature" branch
    And I have an uncommitted file
    When I run `git sync --all` in the "new_folder" folder


  @finishes-with-non-empty-stash
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
    And I am in the project root folder
    And I get the error "Automatic merge failed"
    And I am still on the "current_feature" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH          | COMMAND                           |
      | current_feature | git merge --abort                 |
      |                 | git checkout main                 |
      | main            | git checkout current_feature      |
      | current_feature | git stash pop                     |
      |                 | cd <%= git_folder "new_folder" %> |
    And I am still on the "current_feature" branch
    And I again have my uncommitted file
    And there is no merge in progress
    And I am left with my original commits


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "current_feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH          | COMMAND                                  |
      | current_feature | git commit --no-edit                     |
      |                 | git push                                 |
      |                 | git checkout other-feature               |
      | other-feature   | git merge --no-edit origin/other-feature |
      |                 | git merge --no-edit main                 |
      |                 | git push                                 |
      |                 | git checkout current_feature             |
      | current_feature | git stash pop                            |
      |                 | cd <%= git_folder "new_folder" %>        |
    And I am still on the "current_feature" branch
    And I again have my uncommitted file
    And there is no merge in progress
    And now I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                                  | FILE NAME        |
      | main            | local and remote | conflicting main commit                  | conflicting_file |
      | current_feature | local and remote | conflicting feature commit               | conflicting_file |
      |                 |                  | folder commit                            | new_folder/file1 |
      |                 |                  | conflicting main commit                  | conflicting_file |
      |                 |                  | Merge branch 'main' into current_feature |                  |
      | other-feature   | local and remote | other feature commit                     | file2            |
      |                 |                  | conflicting main commit                  | conflicting_file |
      |                 |                  | Merge branch 'main' into other-feature   |                  |
    And I still have the following committed files
      | BRANCH          | NAME             | CONTENT          |
      | main            | conflicting_file | main content     |
      | current_feature | conflicting_file | resolved content |
      | current_feature | new_folder/file1 |                  |
      | other-feature   | conflicting_file | main content     |
      | other-feature   | file2            |                  |
