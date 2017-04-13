Feature: git-repo when origin is unsupported

  Background:
    When I run `gt repo`


  Scenario: result
    Then I get the error "unsupported hosting service"
    And I get the error "This command requires hosting on GitHub, GitLab, or Bitbucket"
