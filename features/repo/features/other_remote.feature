Feature: display the repo at another remote

  Scenario: upstream remote
    Given a Git repo with origin
    And an additional "upstream" remote with URL "https://github.com/upstream/repo.git"
    And tool "open" is installed
    When I run "git-town repo upstream"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://github.com/upstream/repo
      """
