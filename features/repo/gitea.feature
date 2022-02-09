@skipWindows
Feature: Gitea

  Scenario Outline:
    Given the origin is "<ORIGIN>"
    And the "open" tool is installed
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://gitea.com/git-town/git-town
      """

    Examples:
      | ORIGIN                                    |
      | http://gitea.com/git-town/git-town.git    |
      | http://gitea.com/git-town/git-town        |
      | https://gitea.com/git-town/git-town.git   |
      | https://gitea.com/git-town/git-town       |
      | git@gitea.com:git-town/git-town.git       |
      | git@gitea.com:git-town/git-town           |
      | ssh://git@gitea.com/git-town/git-town.git |
      | ssh://git@gitea.com/git-town/git-town     |
