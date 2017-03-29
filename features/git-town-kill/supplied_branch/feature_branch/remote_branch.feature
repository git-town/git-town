Feature: git town-kill: killing a remote only branch

  Background:
    Given I have a feature branch named "feature" on another machine
    And the following commit exists in my repository on another machine
      | BRANCH  | LOCATION         | MESSAGE        |
      | feature | local and remote | feature commit |
    And I am on the "main" branch
    When I run `gt kill feature`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune        |
      |        | git push origin :feature |
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |


  Scenario: undoing the kill
    When I run `gt kill --undo`
    Then it runs the commands
      | BRANCH | COMMAND                                                        |
      | main   | git push origin <%= sha 'feature commit' %>:refs/heads/feature |
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main          |
      | remote     | main, feature |
