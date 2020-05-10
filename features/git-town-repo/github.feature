Feature: git-repo when origin is on GitHub

  Scenario Outline: result
    Given my repo's origin is "<ORIGIN>"
    And my computer has the "open" tool installed
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town
      """

    Examples:
      | ORIGIN                                     |
      | http://github.com/git-town/git-town.git    |
      | http://github.com/git-town/git-town        |
      | https://github.com/git-town/git-town.git   |
      | https://github.com/git-town/git-town       |
      | git@github.com:git-town/git-town.git       |
      | git@github.com:git-town/git-town           |
      | ssh://git@github.com/git-town/git-town.git |
      | ssh://git@github.com/git-town/git-town     |
