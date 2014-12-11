Feature: git-pr when origin is on GitHub over SSH

  As a developer having finished a feature on a repository hosted on GitHub over SSH
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Background:
    Given I have a feature branch named "feature"
    And my remote origin is on GitHub through SSH
    And I am on the "feature" branch
    When I run `git pr`


  Scenario:
    Then I see a browser window for a new pull request on GitHub for the "feature" branch
