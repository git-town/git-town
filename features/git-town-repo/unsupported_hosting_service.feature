Feature: git-repo when origin is unsupported

  Background:
    When I run `git town-repo`


  Scenario: result
    Then I get the error "Unsupported hosting service"
    And I get the error "This command requires hosting on GitHub or Bitbucket"
