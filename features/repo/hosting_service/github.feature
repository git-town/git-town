@skipWindows
Feature: GitHub

  Scenario Outline:
    Given a Git repo with origin
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                   |
      | main   | frontend | open https://github.com/git-town/git-town |

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
