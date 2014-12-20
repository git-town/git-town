Feature: git-repo when origin is on GitHub

  Scenario Outline: result
    Given my remote origin is on GitHub through <protocol>
    And I have "open" installed
    When I run `git repo`
    Then I see a browser open to the homepage of my GitHub repository

    Examples:
      | protocol                   |
      | HTTP ending with .git      |
      | HTTP not ending with .git  |
      | HTTPS ending with .git     |
      | HTTPS not ending with .git |
      | SSH ending with .git       |
      | SSH not ending with .git   |
