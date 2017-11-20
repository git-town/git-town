Feature: git-repo when origin is on GitLab

  Scenario Outline: result
    Given my repo's remote origin is <ORIGIN>
    And I have "open" installed
    When I run `git-town repo`
    Then I see the GitLab homepage of the "kadu/kadu" repository in my browser

    Examples:
      | ORIGIN                           |
      | https://gitlab.com/kadu/kadu.git |
      | git@gitlab.com:kadu/kadu.git     |
