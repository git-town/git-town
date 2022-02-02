Feature: Bitbucket support

  @skipWindows
  Scenario Outline: normal origin
    Given my repo has a feature branch "feature"
    And my repo's origin is "<ORIGIN>"
    And my computer has the "open" tool installed
    And I am on the "feature" branch
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

  @skipWindows
  Scenario Outline: origin includes path that looks like a URL
    Given my repo has a feature branch "feature"
    And my repo's origin is "<ORIGIN>"
    And my computer has the "open" tool installed
    And I am on the "feature" branch
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

  @skipWindows
  Scenario Outline: SSH style origin
    Given my repo has a feature branch "feature"
    And my repo's origin is "<ORIGIN>"
    And my computer has the "open" tool installed
    And I am on the "feature" branch
    When I run "git-town new-pull-request"
    Then "open" launches a new pull request with this url in my browser:
      """
      https://bitbucket.org/git-town/git-town/pull-request/new?dest=git-town%2Fgit-town%3A%3Amain&source=git-town%2Fgit-town%.*%3Afeature
      """

    Examples:
      | ORIGIN                                        |
      | ssh://git@bitbucket.org/git-town/git-town.git |
      | ssh://git@bitbucket.org/git-town/git-town     |
