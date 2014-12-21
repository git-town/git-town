Feature: git-pr: when origin is unsupported

  As a developer trying to create a pull request for a repository on an unsupported hosting service
  I should get an error that my hosting service is not supported
  So I know why the command does not work.


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `git pr` while allowing errors


  Scenario: result
    Then I get the error "Unsupported hosting service. Pull requests can only be created on Bitbucket and GitHub"
