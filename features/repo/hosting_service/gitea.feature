@skipWindows
Feature: Gitea

  Scenario Outline:
    Given a Git repo with origin
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                  |
      | (none) | frontend | open https://gitea.com/git-town/git-town |

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
