Feature: git-repo when origin is unsupported

  Background:
    When I run `git repo`


  Scenario: result
    Then I get the error "Unsupported hosting service. Repositories can only be viewed on Bitbucket and GitHub"
