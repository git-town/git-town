@skipWindows
Feature: self hosted servie

  Scenario Outline:
    Given the "open" tool is installed
    And the origin is "git@self-hosted:git-town/git-town.git"
    And the "code-hosting-driver" setting is "<DRIVER>"
    When I run "git-town repo"
    Then "open" launches a new pull request with this url in my browser:
      """
      <REPO_URL>
      """

    Examples:
      | DRIVER    | REPO_URL                              |
      | bitbucket | https://self-hosted/git-town/git-town |
      | github    | https://self-hosted/git-town/git-town |
      | gitea     | https://self-hosted/git-town/git-town |
      | gitlab    | https://self-hosted/git-town/git-town |
