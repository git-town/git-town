@skipWindows
Feature: Bitbucket support

  Scenario Outline: normal origin
    Given the current branch is a feature branch "feature"
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://bitbucket.org/git-town/git-town/pull-request/new?dest=git-town%2Fgit-town%3A%3Amain&source=git-town%2Fgit-town%.*%3Afeature
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
    Given the current branch is a feature branch "feature"
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://bitbucket.org/git-town/git-town.github.com/pull-request/new?dest=git-town%2Fgit-town.github.com%3A%3Amain&source=git-town%2Fgit-town.github.com%.*%3Afeature
      """

    Examples:
      | ORIGIN                                                          |
      | http://username@bitbucket.org/git-town/git-town.github.com.git  |
      | http://username@bitbucket.org/git-town/git-town.github.com      |
      | https://username@bitbucket.org/git-town/git-town.github.com.git |
      | https://username@bitbucket.org/git-town/git-town.github.com     |
      | git@bitbucket.org/git-town/git-town.github.com.git              |
      | git@bitbucket.org/git-town/git-town.github.com                  |
