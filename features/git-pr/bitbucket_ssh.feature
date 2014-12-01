Feature: git-pr when origin is on Bitbucket over SSH

  Background:
    Given I have a feature branch named "feature"
    And my remote origin is on Bitbucket through SSH
    And I am on the "feature" branch
    When I run `git pr`


  Scenario:
    Then my browser is opened to a new pull request on Bitbucket for the "feature" branch
