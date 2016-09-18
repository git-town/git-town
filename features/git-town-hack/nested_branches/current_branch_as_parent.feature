Feature: Creating nested feature branches

  As a developer waiting for permission to ship a feature branch that contains changes needed for the next feature
  I want to be able to start working on the next feature while having access to the changes currently under review
  So that I am not slowed down by reviews and can keep working on my backlog.


  Background:
    Given I have a feature branch named "parent-feature"
    And the following commits exist in my repository
      | BRANCH         | LOCATION         | MESSAGE        |
      | parent-feature | local and remote | feature_commit |
    And I am on the "parent-feature" branch
    And I have an uncommitted file


  Scenario: Providing the name of the current branch
    When I run `git town-hack child-feature parent-feature`
    Then it runs the commands
      | BRANCH         | COMMAND                                      |
      | parent-feature | git fetch --prune                            |
      |                | git stash -u                                 |
      |                | git checkout main                            |
      | main           | git rebase origin/main                       |
      |                | git checkout parent-feature                  |
      | parent-feature | git merge --no-edit origin/parent-feature    |
      |                | git merge --no-edit main                     |
      |                | git checkout -b child-feature parent-feature |
      | child-feature  | git push -u origin child-feature             |
      |                | git stash pop                                |
    And I end up on the "child-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE        |
      | child-feature  | local and remote | feature_commit |
      | parent-feature | local and remote | feature_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |


  Scenario: Undo
    Given I run `git town-hack child-feature parent-feature`
    When I run `git town-hack --undo`
    Then it runs the commands
      | BRANCH         | COMMAND                        |
      | child-feature  | git stash -u                   |
      |                | git push origin :child-feature |
      |                | git checkout parent-feature    |
      | parent-feature | git branch -d child-feature    |
      |                | git checkout main              |
      | main           | git checkout parent-feature    |
      | parent-feature | git stash pop                  |
    And I end up on the "parent-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | parent-feature | main   |


  Scenario: Providing '.' as the parent name
    When I run `git town-hack child-feature .`
    Then it runs the commands
      | BRANCH         | COMMAND                                      |
      | parent-feature | git fetch --prune                            |
      |                | git stash -u                                 |
      |                | git checkout main                            |
      | main           | git rebase origin/main                       |
      |                | git checkout parent-feature                  |
      | parent-feature | git merge --no-edit origin/parent-feature    |
      |                | git merge --no-edit main                     |
      |                | git checkout -b child-feature parent-feature |
      | child-feature  | git push -u origin child-feature             |
      |                | git stash pop                                |
    And I end up on the "child-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH         | LOCATION         | MESSAGE        |
      | child-feature  | local and remote | feature_commit |
      | parent-feature | local and remote | feature_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT         |
      | child-feature  | parent-feature |
      | parent-feature | main           |


  Scenario: Undo
    Given I run `git town-hack child-feature .`
    When I run `git town-hack --undo`
    Then it runs the commands
      | BRANCH         | COMMAND                        |
      | child-feature  | git stash -u                   |
      |                | git push origin :child-feature |
      |                | git checkout parent-feature    |
      | parent-feature | git branch -d child-feature    |
      |                | git checkout main              |
      | main           | git checkout parent-feature    |
      | parent-feature | git stash pop                  |
    And I end up on the "parent-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH         | PARENT |
      | parent-feature | main   |
