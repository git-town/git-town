Feature: git ship: shipping a parent branch

  (see ../../current_branch/on_feature_branch/on_parent_branch.feature)


  Background:
    Given I have a feature branch named "parent-feature"
    And I have a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION         | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | parent-feature | local and remote | parent feature commit | parent_feature_file | parent feature content |
      | child-feature  | local and remote | child feature commit  | child_feature_file  | child feature content  |
    And I am on the "child-feature" branch
    When I run `git ship parent-feature -m "parent feature done"`


  Scenario: result
    Then it runs the Git commands
      | BRANCH         | COMMAND                                   |
      | child-feature  | git fetch --prune                         |
      |                | git checkout main                         |
      | main           | git rebase origin/main                    |
      |                | git checkout parent-feature               |
      | parent-feature | git merge --no-edit origin/parent-feature |
      |                | git merge --no-edit main                  |
      |                | git checkout main                         |
      | main           | git merge --squash parent-feature         |
      |                | git commit -m "parent feature done"       |
      |                | git push                                  |
      |                | git push origin :parent-feature           |
      |                | git branch -D parent-feature              |
      |                | git checkout child-feature                |
    And I end up on the "child-feature" branch
    And I have the following commits
      | BRANCH        | LOCATION         | MESSAGE              | FILE NAME           | FILE CONTENT           |
      | main          | local and remote | parent feature done  | parent_feature_file | parent feature content |
      | child-feature | local and remote | child feature commit | child_feature_file  | child feature content  |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | child-feature | main   |
