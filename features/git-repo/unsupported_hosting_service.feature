Feature: git-repo when origin is unsupported

  Background:
    When I run `git repo` while allowing errors


  Scenario:
    Then I get the error "Unsupported hosting service. Repositories can only be opened from Bitbucket and GitHub"
