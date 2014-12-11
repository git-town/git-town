Feature: git-pr when origin is on Bitbucket over SSH

  As a developer having finished a feature on a repository hosted on Bitbucket over SSH
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Background:
    Given I have a feature branch named "feature"
    And my remote origin is on Bitbucket through SSH
    And I am on the "feature" branch
    When I run `git pr`


  Scenario:
    Then I see a browser window for a new pull request on Bitbucket for the "feature" branch
