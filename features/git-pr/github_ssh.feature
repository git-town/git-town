Feature: git-pr when origin is on GitHub over SSH

  Background:
    Given I have a feature branch named "feature"
    And my remote origin is on GitHub through SSH
    And I am on the "feature" branch
    When I run `git pr`


  Scenario: result
    Then I see a browser window for a new pull request on GitHub for the "feature" branch
