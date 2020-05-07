Feature: git town-kill: killing a remote only branch

  Background:
    Given my origin has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | remote   | feature commit |
    And I am on the "main" branch
    And I run "git-town sync"
    When I run "git-town kill feature"


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git push origin :feature |
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |


  Scenario: undoing the kill
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                 |
      | main   | git push origin {{ sha-in-remote 'feature commit' }}:refs/heads/feature |
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main          |
      | remote     | main, feature |
