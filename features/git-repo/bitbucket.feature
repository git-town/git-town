Feature: git-repo when origin is on Bitbucket

  Scenario Outline: result
    Given my remote origin is <REPOSITORY>
    And I have "open" installed
    When I run `git repo`
    Then I see the homepage of my Bitbucket repository in my browser

    Examples:
      | REPOSITORY                                            |
      | http://username@bitbucket.org/Originate/git-town.git  |
      | http://username@bitbucket.org/Originate/git-town      |
      | https://username@bitbucket.org/Originate/git-town.git |
      | https://username@bitbucket.org/Originate/git-town     |
      | git@bitbucket.org/Originate/git-town.git              |
      | git@bitbucket.org/Originate/git-town                  |
