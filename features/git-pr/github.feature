Feature: git-pr when origin is on GitHub

  As a developer having finished a feature in a repository hosted on GitHub
  I want to be able to easily create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Scenario Outline: result
    Given I have a feature branch named "feature"
    And my remote origin is on GitHub through <PROTOCOL>
    And I am on the "feature" branch
    When I run `git pr`
    Then I see a browser window for a new pull request on GitHub for the "feature" branch

    Examples:
      | PROTOCOL                   |
      | HTTP ending with .git      |
      | HTTP not ending with .git  |
      | HTTPS ending with .git     |
      | HTTPS not ending with .git |
      | SSH ending with .git       |
      | SSH not ending with .git   |
