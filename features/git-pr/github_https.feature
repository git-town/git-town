Feature: git-pr when origin is on GitHub over HTTPS

  Background:
    Given I have a feature branch named "feature"
    And my remote origin is on GitHub through HTTPS
    And I am on the "feature" branch
    When I run `git pr`


  Scenario:
    Then I see a browser window for a new pull request on GitHub for the "feature" branch
