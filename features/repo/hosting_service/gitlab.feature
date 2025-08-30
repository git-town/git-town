@skipWindows
Feature: GitLab

  Scenario Outline:
    Given a Git repo with origin
    And the origin is "<ORIGIN>"
    And tool "open" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                           |
      | main   | frontend | open https://gitlab.com/kadu/kadu |

    Examples:
      | ORIGIN                           |
      | https://gitlab.com/kadu/kadu.git |
      | git@gitlab.com:kadu/kadu.git     |
