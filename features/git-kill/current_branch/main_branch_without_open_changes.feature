Feature: Git Kill: The main branch is not killed

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
