Feature: git-new-pull-request: when origin is unsupported

  As a developer trying to create a pull request in a repository on an unsupported hosting service
  I should get an error that my hosting service is not supported
  So that I know why the command doesn't work.


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `gt new-pull-request`


  Scenario: result
    Then I get the error "unsupported hosting service"
    And I get the error "This command requires hosting on GitHub, GitLab, or Bitbucket"
