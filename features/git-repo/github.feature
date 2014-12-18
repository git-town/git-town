Feature: git-repo when origin is on GitHub

  Scenario Outline: result
    Given my remote origin is on GitHub through <PROTOCOL>
    When I run `git repo`
    Then I see a browser window for my repository homepage on GitHub

    Examples:
      | PROTOCOL                   |
      | HTTP ending with .git      |
      | HTTP not ending with .git  |
      | HTTPS ending with .git     |
      | HTTPS not ending with .git |
      | SSH ending with .git       |
      | SSH not ending with .git   |
