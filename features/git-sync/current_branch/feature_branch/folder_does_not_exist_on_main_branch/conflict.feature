Feature: git sync: syncing inside a folder that doesn't exist on the main branch (with open changes)

  (see ./no_conflict.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local and remote | conflicting main commit    | conflicting_file | main content    |
      | feature | local            | conflicting feature commit | conflicting_file | feature content |
      |         |                  | folder commit              | new_folder/file1 |                 |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git sync` in the "new_folder" folder


  @finishes-with-non-empty-stash
  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      |         | cd <%= git_root_folder %>          |
      |         | git stash -u                       |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And I am in the project root folder
    And I get the error "Automatic merge failed"
    And I am still on the "feature" branch
    And I don't have my uncommitted file
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH  | COMMAND                           |
      | feature | git merge --abort                 |
      |         | git checkout main                 |
      | main    | git checkout feature              |
      | feature | git stash pop                     |
      |         | cd <%= git_folder "new_folder" %> |
    And I am still on the "feature" branch
    And I again have my uncommitted file
    And there is no merge in progress
    And I am left with my original commits


  @finishes-with-non-empty-stash
  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature" branch
    And I don't have my uncommitted file
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH  | COMMAND                           |
      | feature | git commit --no-edit              |
      |         | git push                          |
      |         | git stash pop                     |
      |         | cd <%= git_folder "new_folder" %> |
    And I am still on the "feature" branch
    And I again have my uncommitted file
    And there is no merge in progress
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                          | FILE NAME        |
      | main    | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | conflicting feature commit       | conflicting_file |
      |         |                  | folder commit                    | new_folder/file1 |
      |         |                  | conflicting main commit          | conflicting_file |
      |         |                  | Merge branch 'main' into feature |                  |
    And I still have the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
      | feature | new_folder/file1 |                  |
