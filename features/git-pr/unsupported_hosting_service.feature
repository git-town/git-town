Feature: git-pr when origin is unsupported

  As a developer trying to create a pull request for a repository on an unknown hosting service
  I want to get a clear error message explaining why the feature doesn't work, and what to do to make it work
  So that I can configure and use this tool the right way without having to read the manual and have more time for coding.


  Background:
    Given I have a feature branch named "feature"
    And I am on the "feature" branch
    When I run `git pr` while allowing errors


  Scenario:
    Then I get the error "Unsupported hosting service. Pull requests can only be created on Bitbucket and GitHub"
