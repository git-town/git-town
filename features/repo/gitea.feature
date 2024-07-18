@skipWindows
Feature: Gitea

  Scenario Outline:
    Given a Git repo clone
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
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
