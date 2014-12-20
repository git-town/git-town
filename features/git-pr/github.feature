Feature: git-pr when origin is on GitHub

  Scenario Outline: result
    Given I have a feature branch named "feature"
    And my remote origin is on GitHub through <protocol>
    And I have "open" installed
    And I am on the "feature" branch
    When I run `git pr`
    Then I see a browser open to a new GitHub pull request for the "feature" branch

    Examples:
      | protocol                   |
      | HTTP ending with .git      |
      | HTTP not ending with .git  |
      | HTTPS ending with .git     |
      | HTTPS not ending with .git |
      | SSH ending with .git       |
      | SSH not ending with .git   |
