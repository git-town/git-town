Feature: shipping a parent branch

  Background:
    Given my repo has a feature branch named "parent-feature"
    And my repo has a feature branch named "child-feature" as a child of "parent-feature"
    And the following commits exist in my repo
      | BRANCH         | LOCATION      | MESSAGE               |
      | parent-feature | local, remote | parent feature commit |
      | child-feature  | local, remote | child feature commit  |
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
      | BRANCH         | LOCATION      | MESSAGE               |
      | main           | local, remote | parent feature done   |
      | child-feature  | local, remote | child feature commit  |
      | parent-feature | remote        | parent feature commit |
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
      | BRANCH         | LOCATION      | MESSAGE                      |
      | main           | local, remote | parent feature done          |
      |                |               | Revert "parent feature done" |
      | child-feature  | local, remote | child feature commit         |
      | parent-feature | local, remote | parent feature commit        |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |
