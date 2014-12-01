Feature: git-pr when origin is on Bitbucket over HTTPS

  Background:
    Given I have a feature branch named "feature"
    And my remote origin is on Bitbucket through HTTPS
    And I am on the "feature" branch
    When I run `git pr`


  Scenario:
    Then I see a browser open to a new pull request on Bitbucket for the "feature" branch
