Feature: Skip deleting the remote branch when shipping

  When using GitHub's feature to automatically delete head branches of pull requests.
  I want "git ship" to skip deleting the remote feature branch
  So that I can keep using Git Town in this situation.


  Background:
    Given my repository has a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local and remote | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    And I have a the git configuration for "git-town.ship-delete-remote-branch" set to "false"
    When I run `git-town ship -m "feature done"`


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
      |         | git branch -D feature              |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main          |
      | remote     | main, feature |
    And my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    |
      | main    | local and remote | feature done   | feature_file |
      | feature | remote           | feature commit | feature_file |


  Scenario: undo
    When I run `git-town undo`
    Then it runs the commands
      | BRANCH  | COMMAND                                        |
      | main    | git branch feature <%= sha 'feature commit' %> |
      |         | git revert <%= sha 'feature done' %>           |
      |         | git push                                       |
      |         | git checkout feature                           |
      | feature | git checkout main                              |
      | main    | git checkout feature                           |
    And I end up on the "feature" branch
    And my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE               | FILE NAME    |
      | main    | local and remote | feature done          | feature_file |
      |         |                  | Revert "feature done" | feature_file |
      | feature | local and remote | feature commit        | feature_file |
