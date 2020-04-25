Feature: git-repo when origin is on Bitbucket

  Scenario Outline: result
    Given my repo's remote origin is <ORIGIN>
    And I have "open" installed
    When I run `git-town repo`
    Then I see the Bitbucket homepage of the "git-town/git-town" repository in my browser

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
