@skipWindows
Feature: GitHub

  Scenario Outline:
    Given the origin is "<ORIGIN>"
    And the "open" tool is installed
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
