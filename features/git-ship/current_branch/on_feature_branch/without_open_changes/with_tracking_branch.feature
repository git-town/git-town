Feature: git ship: shipping the current feature branch with a tracking branch

  As a developer having finished a feature
  I want to be able to ship it safely in one easy step
  So that I can quickly move on to the next feature and remain productive.


  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local and remote | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run `git ship -m "feature done"`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
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
    And I end up on the "main" branch
    And there are no more feature branches
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |


  Scenario: undo
    When I run `git ship --undo`
    Then it runs the Git commands
      | BRANCH  | COMMAND                                        |
      | main    | git branch feature <%= sha 'feature commit' %> |
      |         | git push -u origin feature                     |
      |         | git revert <%= sha 'feature done' %>           |
      |         | git push origin main                           |
      |         | git checkout feature                           |
      | feature | git checkout main                              |
      | main    | git checkout feature                           |
    And I end up on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE               | FILE NAME    |
      | main    | local and remote | feature done          | feature_file |
      |         |                  | Revert "feature done" | feature_file |
      | feature | local and remote | feature commit        | feature_file |

