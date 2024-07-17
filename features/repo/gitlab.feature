@skipWindows
Feature: GitLab

  Scenario Outline:
    Given a Git repo clone
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://gitlab.com/kadu/kadu
      """

    Examples:
      | ORIGIN                           |
      | https://gitlab.com/kadu/kadu.git |
      | git@gitlab.com:kadu/kadu.git     |
