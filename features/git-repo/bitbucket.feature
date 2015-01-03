Feature: git-repo when origin is on Bitbucket

  Scenario Outline: result
    Given my remote origin is <ORIGIN>
    And I have "open" installed
    When I run `git repo`
    Then I see the Bitbucket homepage of the "<REPOSITORY>" repository in my browser

    Examples:
      | ORIGIN                                                | REPOSITORY         |
      | http://username@bitbucket.org/Originate/git-town.git  | Originate/git-town |
      | http://username@bitbucket.org/Originate/git-town      | Originate/git-town |
      | https://username@bitbucket.org/Originate/git-town.git | Originate/git-town |
      | https://username@bitbucket.org/Originate/git-town     | Originate/git-town |
      | git@bitbucket.org/Originate/git-town.git              | Originate/git-town |
      | git@bitbucket.org/Originate/git-town                  | Originate/git-town |
