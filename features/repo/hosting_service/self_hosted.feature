@skipWindows
Feature: self hosted servie

  Background:
    Given a Git repo with origin

  Scenario Outline:
    Given the origin is "git@self-hosted:git-town/git-town.git"
    And Git setting "git-town.forge-type" is "<DRIVER>"
    And tool "open" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                                    |
      | main   | frontend | open https://self-hosted/git-town/git-town |

    Examples:
      | DRIVER    |
      | bitbucket |
      | github    |
      | gitlab    |
# uncomment to test (makes online connection)
# | gitea     |

  Scenario: GitLab with custom port
    Given the origin is "ssh://git@git.example.com:4022/a/b.git"
    And Git setting "git-town.forge-type" is "gitlab"
    And tool "open" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND                          |
      | main   | frontend | open https://git.example.com/a/b |
