@skipWindows
Feature: Bitbucket support

  Background:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"

  Scenario Outline: normal origin
    Given the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://bitbucket.org/git-town/git-town/pull-requests/new?source=feature&dest=git-town%2Fgit-town%3Amain
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

  Scenario Outline: origin includes path that looks like a URL
    Given the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town propose"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://bitbucket.org/git-town/git-town.github.com/pull-requests/new?source=feature&dest=git-town%2Fgit-town.github.com%3Amain
      """

    Examples:
      | ORIGIN                                                          |
      | http://username@bitbucket.org/git-town/git-town.github.com.git  |
      | http://username@bitbucket.org/git-town/git-town.github.com      |
      | https://username@bitbucket.org/git-town/git-town.github.com.git |
      | https://username@bitbucket.org/git-town/git-town.github.com     |
      | git@bitbucket.org/git-town/git-town.github.com.git              |
      | git@bitbucket.org/git-town/git-town.github.com                  |
