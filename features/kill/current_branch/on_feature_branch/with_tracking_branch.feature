Feature: git town-kill: killing the current feature branch with a tracking branch

  As a developer currently working on a feature branch that leads nowhere
  I want to be able to cleanly delete the current branch including all open changes
  So that my workspace doesn't contain irrelevant branches and my productivity remains high.


  Background:
    Given I have feature branches named "current-feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                |
      | current-feature | local and remote | current feature commit |
      | other-feature   | local and remote | other feature commit   |
    And I am on the "current-feature" branch
    And I have an uncommitted file
    When I run `git-town kill`


  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                                |
      | current-feature | git fetch --prune                      |
      |                 | git push origin :current-feature       |
      |                 | git add -A                             |
      |                 | git commit -m "WIP on current-feature" |
      |                 | git checkout main                      |
      | main            | git branch -D current-feature          |
    And I end up on the "main" branch
    And I don't have any uncommitted files
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And I have the following commits
      | BRANCH        | LOCATION         | MESSAGE              |
      | other-feature | local and remote | other feature commit |


  Scenario: undoing the kill
    When I run `git-town kill --undo`
    Then it runs the commands
      | BRANCH          | COMMAND                                                        |
      | main            | git branch current-feature <%= sha 'WIP on current-feature' %> |
      |                 | git checkout current-feature                                   |
      | current-feature | git reset <%= sha 'current feature commit' %>                  |
      |                 | git push -u origin current-feature                             |
    And I end up on the "current-feature" branch
    And I again have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES                             |
      | local      | main, current-feature, other-feature |
      | remote     | main, current-feature, other-feature |
    And I am left with my original commits
