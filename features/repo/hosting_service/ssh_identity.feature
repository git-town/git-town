@skipWindows
Feature: use an SSH identity

  Scenario Outline:
    Given a Git repo with origin
    And the origin is "git@my-ssh-identity:git-town/git-town.git"
    And Git setting "git-town.hosting-origin-hostname" is "<ORIGIN_HOSTNAME>"
    And tool "open" is installed
    When I run "git-town repo"
    Then Git Town runs the commands
      | BRANCH | TYPE     | COMMAND         |
      | main   | frontend | open <REPO_URL> |

    Examples:
      | ORIGIN_HOSTNAME | REPO_URL                                |
      | bitbucket.org   | https://bitbucket.org/git-town/git-town |
      | github.com      | https://github.com/git-town/git-town    |
      | gitlab.com      | https://gitlab.com/git-town/git-town    |
  # uncomment to test (makes online connection)
  # | gitea.com       | https://gitea.com/git-town/git-town     |
