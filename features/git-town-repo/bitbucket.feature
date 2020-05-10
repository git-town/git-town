Feature: git-repo when origin is on Bitbucket

  Scenario Outline: result
    Given my repo's origin is "<ORIGIN>"
    And my computer has the "open" tool installed
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
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
