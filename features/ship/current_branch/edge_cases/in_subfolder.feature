Feature: ship the current feature branch from a subfolder on the shipped branch

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME               | FILE CONTENT    |
      | feature | local, remote | feature commit | new_folder/feature_file | feature content |
    And I am on the "feature" branch
    When I run "git-town ship -m 'feature done'" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME               |
      | main   | local, remote | feature done | new_folder/feature_file |
    And Git Town now has no branch hierarchy information

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git push -u origin feature                    |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git push                                      |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git checkout feature                          |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE               | FILE NAME               |
      | main    | local, remote | feature done          | new_folder/feature_file |
      |         |               | Revert "feature done" | new_folder/feature_file |
      | feature | local, remote | feature commit        | new_folder/feature_file |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
