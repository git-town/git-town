Feature: git town-kill: killing the current feature branch without a tracking branch

  As a developer currently working on a local feature branch that leads nowhere
  I want to be able to remove the current branch including all open changes
  So that my workspace doesn't contain irrelevant branches and my productivity remains high.


  Background:
    Given I have a feature branch named "other-feature"
    And I have a local feature branch named "current-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | current-feature | local            | current feature commit |
      | other-feature   | local and remote | other feature commit   |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `gt kill`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                                |
      | current-feature | git fetch --prune                      |
      |                 | git add -A                             |
      |                 | git commit -m 'WIP on current-feature' |
      |                 | git checkout main                      |
      | main            | git branch -D current-feature          |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And I have the following commits
      | BRANCH        | LOCATION         | MESSAGE              |
      | other-feature | local and remote | other feature commit |


  Scenario: Undoing a kill of a local feature branch
    When I run `gt kill --undo`
    Then it runs the commands
      | BRANCH          | COMMAND                                                        |
      | main            | git branch current-feature <%= sha 'WIP on current-feature' %> |
      |                 | git checkout current-feature                                   |
      | current-feature | git reset <%= sha 'current feature commit' %>                  |
    And I end up on the "current-feature" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
      | remote     | main, other-feature                  |
    And I am left with my original commits
