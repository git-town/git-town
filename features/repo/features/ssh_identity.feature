@skipWindows
Feature: use an SSH identity

  Scenario Outline:
    Given a Git repo clone
    And tool "open" is installed
    And the origin is "git@my-ssh-identity:git-town/git-town.git"
    And Git Town setting "hosting-origin-hostname" is "<ORIGIN_HOSTNAME>"
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
      """
      <REPO_URL>
      """

    Examples:
      | ORIGIN_HOSTNAME | REPO_URL                                |
      | bitbucket.org   | https://bitbucket.org/git-town/git-town |
      | github.com      | https://github.com/git-town/git-town    |
      | gitea.com       | https://gitea.com/git-town/git-town     |
      | gitlab.com      | https://gitlab.com/git-town/git-town    |
