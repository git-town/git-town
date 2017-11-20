Feature: git-repo when origin is on Bitbucket

  Scenario Outline: result
    Given my repo's remote origin is <ORIGIN>
    And I have "open" installed
    When I run `git-town repo`
    Then I see the Bitbucket homepage of the "Originate/git-town" repository in my browser

    Examples:
      | ORIGIN                                                |
      | http://username@bitbucket.org/Originate/git-town.git  |
      | http://username@bitbucket.org/Originate/git-town      |
      | https://username@bitbucket.org/Originate/git-town.git |
      | https://username@bitbucket.org/Originate/git-town     |
      | git@bitbucket.org/Originate/git-town.git              |
      | git@bitbucket.org/Originate/git-town                  |
      | ssh://git@bitbucket.org/Originate/git-town.git        |
      | ssh://git@bitbucket.org/Originate/git-town            |
