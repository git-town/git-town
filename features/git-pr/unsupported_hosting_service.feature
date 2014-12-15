Feature: git-pr when origin is unsupported

  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `git pr` while allowing errors


  Scenario: result
    Then I get the error "Unsupported hosting service. Pull requests can only be created on Bitbucket and GitHub"
