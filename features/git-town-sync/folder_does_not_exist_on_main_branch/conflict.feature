Feature: git-town sync: syncing inside a folder that doesn't exist on the main branch

  (see ./no_conflict.feature)


  Background:
    Given I have feature branches named "current-feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main            | local and remote | conflicting main commit    | conflicting_file | main content    |
      | current-feature | local            | conflicting feature commit | conflicting_file | feature content |
      |                 |                  | folder commit              | new_folder/file1 |                 |
      | other-feature   | local and remote | other feature commit       | file2            |                 |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `git-town sync --all` in the "new_folder" folder


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                                    |
      | current-feature | git fetch --prune                          |
      | <none>          | cd <%= git_root_folder %>                  |
      | current-feature | git add -A                                 |
      |                 | git stash                                  |
      |                 | git checkout main                          |
      | main            | git rebase origin/main                     |
      |                 | git checkout current-feature               |
      | current-feature | git merge --no-edit origin/current-feature |
      |                 | git merge --no-edit main                   |
    And I am in the project root folder
    And I get the error "Automatic merge failed"
    And I am still on the "current-feature" branch
    And my uncommitted file is stashed
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git-town sync --abort`
    Then it runs the commands
      | BRANCH          | COMMAND                           |
      | current-feature | git merge --abort                 |
      |                 | git checkout main                 |
      | main            | git checkout current-feature      |
      | current-feature | git stash pop                     |
      | <none>          | cd <%= git_folder "new_folder" %> |
    And I am still on the "current-feature" branch
    And I again have my uncommitted file
    And there is no merge in progress
    And I am left with my original commits


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving the conflicts
    When I run `git-town sync --continue`
    Then it runs no commands
    And I get the error "You must resolve the conflicts before continuing"
    And I am still on the "current-feature" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git-town sync --continue`
    Then it runs the commands
      | BRANCH          | COMMAND                                  |
      | current-feature | git commit --no-edit                     |
      |                 | git push                                 |
      |                 | git checkout other-feature               |
      | other-feature   | git merge --no-edit origin/other-feature |
      |                 | git merge --no-edit main                 |
      |                 | git push                                 |
      |                 | git checkout current-feature             |
      | current-feature | git push --tags                          |
      |                 | git stash pop                            |
      | <none>          | cd <%= git_folder "new_folder" %>        |
    And I am still on the "current-feature" branch
    And I again have my uncommitted file
    And there is no merge in progress
    And now I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                                  | FILE NAME        |
      | main            | local and remote | conflicting main commit                  | conflicting_file |
      | current-feature | local and remote | conflicting feature commit               | conflicting_file |
      |                 |                  | folder commit                            | new_folder/file1 |
      |                 |                  | conflicting main commit                  | conflicting_file |
      |                 |                  | Merge branch 'main' into current-feature |                  |
      | other-feature   | local and remote | other feature commit                     | file2            |
      |                 |                  | conflicting main commit                  | conflicting_file |
      |                 |                  | Merge branch 'main' into other-feature   |                  |
    And I still have the following committed files
      | BRANCH          | NAME             | CONTENT          |
      | main            | conflicting_file | main content     |
      | current-feature | conflicting_file | resolved content |
      | current-feature | new_folder/file1 |                  |
      | other-feature   | conflicting_file | main content     |
      | other-feature   | file2            |                  |
