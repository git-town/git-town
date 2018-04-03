Feature: git town-ship: shipping the current feature branch from a subfolder

  As a developer shipping a feature branch from a subfolder
  I want the command to finish properly
  So that my repo is left in a consistent state and I don't lose any data


  Background:
    Given my repository has a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME               | FILE CONTENT    |
      | feature | local and remote | feature commit | new_folder/feature_file | feature content |
    And I am on the "feature" branch
    When I run `git-town ship -m "feature done"` in the "new_folder" folder


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      | <none>  | cd <%= git_root_folder %>          |
      | feature | git checkout main                  |
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
      | <none>  | cd <%= git_folder "new_folder" %>  |
    And I end up on the "main" branch
    And there are no more feature branches
    And my repository has the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME               |
      | main   | local and remote | feature done | new_folder/feature_file |


  Scenario: undo
    When I run `git-town undo`
    Then it runs the commands
      | BRANCH  | COMMAND                                        |
      | <none>  | cd <%= git_root_folder %>                      |
      | main    | git branch feature <%= sha 'feature commit' %> |
      |         | git push -u origin feature                     |
      |         | git revert <%= sha 'feature done' %>           |
      |         | git push                                       |
      |         | git checkout feature                           |
      | feature | git checkout main                              |
      | main    | git checkout feature                           |
      | <none>  | cd <%= git_folder "new_folder" %>              |
    And I end up on the "feature" branch
    And my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE               | FILE NAME               |
      | main    | local and remote | feature done          | new_folder/feature_file |
      |         |                  | Revert "feature done" | new_folder/feature_file |
      | feature | local and remote | feature commit        | new_folder/feature_file |
