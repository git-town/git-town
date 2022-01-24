Feature: git town-ship: shipping a parent branch


  Background:
    Given my repo has a feature branch named "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repo
      | BRANCH         | LOCATION      | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | parent-feature | local, remote | parent feature commit | parent_feature_file | parent feature content |
      | child-feature  | local, remote | child feature commit  | child_feature_file  | child feature content  |
    And I am on the "parent-feature" branch
    When I run "git-town ship -m 'parent feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH         | COMMAND                                   |
      | parent-feature | git fetch --prune --tags                  |
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
    And I am now on the "main" branch
    And my repo now has the following commits
      | BRANCH         | LOCATION      | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | main           | local, remote | parent feature done   | parent_feature_file | parent feature content |
      | child-feature  | local, remote | child feature commit  | child_feature_file  | child feature content  |
      | parent-feature | remote        | parent feature commit | parent_feature_file | parent feature content |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | child-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH         | COMMAND                                                     |
      | main           | git branch parent-feature {{ sha 'parent feature commit' }} |
      |                | git revert {{ sha 'parent feature done' }}                  |
      |                | git push                                                    |
      |                | git checkout parent-feature                                 |
      | parent-feature | git checkout main                                           |
      | main           | git checkout parent-feature                                 |
    And I am now on the "parent-feature" branch
    And my repo now has the following commits
      | BRANCH         | LOCATION      | MESSAGE                      | FILE NAME           |
      | main           | local, remote | parent feature done          | parent_feature_file |
      |                |               | Revert "parent feature done" | parent_feature_file |
      | child-feature  | local, remote | child feature commit         | child_feature_file  |
      | parent-feature | local, remote | parent feature commit        | parent_feature_file |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |
