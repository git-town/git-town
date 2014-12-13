Feature: git-repo when origin is unsupported

  Background:
    When I run `git repo` while allowing errors


  Scenario: result
    Then I get the error "Unsupported hosting service. Can only view repositories on Bitbucket and GitHub"
