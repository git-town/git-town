@skipWindows
Feature: use an SSH identity

  Scenario Outline:
    Given my computer has the "open" tool installed
    And my repo's origin is "git@my-ssh-identity:git-town/git-town.git"
    And the "code-hosting-origin-hostname" setting is "<ORIGIN_HOSTNAME>"
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      <REPO_URL>
      """

    Examples:
      | ORIGIN_HOSTNAME | REPO_URL                                |
      | bitbucket.org   | https://bitbucket.org/git-town/git-town |
      | github.com      | https://github.com/git-town/git-town    |
      | gitea.com       | https://gitea.com/git-town/git-town     |
      | gitlab.com      | https://gitlab.com/git-town/git-town    |
