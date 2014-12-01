Feature: git-pr when origin is on Github over SSH

  Background:
    Given I have a feature branch named "feature"
    And my remote origin is on Github through SSH
    And I am on the "feature" branch
    When I run `git pr`


  Scenario:
    Then my browser is opened to a new pull request on Github for the "feature" branch
