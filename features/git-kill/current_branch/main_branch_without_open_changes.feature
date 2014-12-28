Feature: git kill: don't remove the main branch (without open changes)

  (see ./main_branch_with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And I am on the "main" branch
    When I run `git kill` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "You can only kill feature branches"
    And I am still on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
