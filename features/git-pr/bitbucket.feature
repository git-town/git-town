Feature: git-pr when origin is on Bitbucket

  Scenario Outline: result
    Given I have a feature branch named "feature"
    And my remote origin is on Bitbucket through <protocol>
    And I am on the "feature" branch
    When I run `git pr`
    Then I see a browser window for a new pull request on Bitbucket for the "feature" branch

    Examples:
      | protocol                   |
      | HTTP ending with .git      |
      | HTTP not ending with .git  |
      | HTTPS ending with .git     |
      | HTTPS not ending with .git |
      | SSH ending with .git       |
      | SSH not ending with .git   |
