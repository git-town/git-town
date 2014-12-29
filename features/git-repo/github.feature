Feature: git-repo when origin is on GitHub

  Scenario Outline: result
    Given my remote origin is <ORIGIN>
    And I have "open" installed
    When I run `git repo`
    Then I see the homepage of my GitHub repository in my browser

    Examples:
      | ORIGIN                                    |
      | http://github.com/Originate/git-town.git  |
      | http://github.com/Originate/git-town      |
      | https://github.com/Originate/git-town.git |
      | https://github.com/Originate/git-town     |
      | git@github.com:Originate/git-town.git     |
      | git@github.com:Originate/git-town         |
