Feature: git-repo when origin is on Bitbucket

  Scenario Outline: result
    Given my remote origin is on Bitbucket through <protocol>
    When I run `git repo`
    Then I see a browser window for my repository homepage on Bitbucket

    Examples:
      | protocol                   |
      | HTTP ending with .git      |
      | HTTP not ending with .git  |
      | HTTPS ending with .git     |
      | HTTPS not ending with .git |
      | SSH ending with .git       |
      | SSH not ending with .git   |
