@skipWindows
Feature: Bitbucket

  Scenario Outline:
    Given a Git repo clone
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://bitbucket.org/git-town/git-town
      """

    Examples:
      | ORIGIN                                               |
      | http://username@bitbucket.org/git-town/git-town.git  |
      | http://username@bitbucket.org/git-town/git-town      |
      | https://username@bitbucket.org/git-town/git-town.git |
      | https://username@bitbucket.org/git-town/git-town     |
      | git@bitbucket.org/git-town/git-town.git              |
      | git@bitbucket.org/git-town/git-town                  |
      | ssh://git@bitbucket.org/git-town/git-town.git        |
      | ssh://git@bitbucket.org/git-town/git-town            |
      | username@bitbucket.org/git-town/git-town.git         |
      | username@bitbucket.org/git-town/git-town             |
      | ssh://username@bitbucket.org/git-town/git-town.git   |
      | ssh://username@bitbucket.org/git-town/git-town       |
