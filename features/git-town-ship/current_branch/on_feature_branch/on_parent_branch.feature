Feature: git town-ship: shipping a parent branch

  As a user shipping a feature branch that is a parent branch to other feature branches
  I want that the child branches are direct descendents of main after shipping
  So that my workspace stays in a consistent state at all times.


  Background:
    Given my repository has a feature branch named "parent-feature"
    And my repository has a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | parent-feature | local and remote | parent feature commit | parent_feature_file | parent feature content |
      | child-feature  | local and remote | child feature commit  | child_feature_file  | child feature content  |
    And I am on the "parent-feature" branch
    When I run `git-town ship -m "parent feature done"`


  Scenario: result
    Then it runs the commands
      | BRANCH         | COMMAND                                   |
      | parent-feature | git fetch --prune                         |
      |                | git checkout main                         |
      | main           | git rebase origin/main                    |
      |                | git checkout parent-feature               |
      | parent-feature | git merge --no-edit origin/parent-feature |
      |                | git merge --no-edit main                  |
      |                | git checkout main                         |
      | main           | git merge --squash parent-feature         |
      |                | git commit -m "parent feature done"       |
      |                | git push                                  |
      |                | git branch -D parent-feature              |
    And I end up on the "main" branch
    And my repository has the following commits
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | main           | local and remote | parent feature done   | parent_feature_file | parent feature content |
      | child-feature  | local and remote | child feature commit  | child_feature_file  | child feature content  |
      | parent-feature | remote           | parent feature commit | parent_feature_file | parent feature content |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | child-feature | main   |


  Scenario: undo
    When I run `git-town undo`
    Then it runs the commands
      | BRANCH         | COMMAND                                                      |
      | main           | git branch parent-feature <%= sha 'parent feature commit' %> |
      |                | git revert <%= sha 'parent feature done' %>                  |
      |                | git push                                                     |
      |                | git checkout parent-feature                                  |
      | parent-feature | git checkout main                                            |
      | main           | git checkout parent-feature                                  |
    And I end up on the "parent-feature" branch
    And my repository has the following commits
      | BRANCH         | LOCATION         | MESSAGE                      | FILE NAME           |
      | main           | local and remote | parent feature done          | parent_feature_file |
      |                |                  | Revert "parent feature done" | parent_feature_file |
      | child-feature  | local and remote | child feature commit         | child_feature_file  |
      | parent-feature | local and remote | parent feature commit        | parent_feature_file |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |
