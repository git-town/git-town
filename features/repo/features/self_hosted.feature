@skipWindows
Feature: self hosted servie

  Background:
    Given a Git repo clone

  Scenario Outline:
    Given tool "open" is installed
    And the origin is "git@self-hosted:git-town/git-town.git"
    And Git Town setting "hosting-platform" is "<DRIVER>"
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
      """
      <REPO_URL>
      """

    Examples:
      | DRIVER    | REPO_URL                              |
      | bitbucket | https://self-hosted/git-town/git-town |
      | github    | https://self-hosted/git-town/git-town |
      | gitea     | https://self-hosted/git-town/git-town |
      | gitlab    | https://self-hosted/git-town/git-town |

  Scenario: GitLab with custom port
    Given the origin is "ssh://git@git.example.com:4022/a/b.git"
    And Git Town setting "hosting-platform" is "gitlab"
    And tool "open" is installed
    When I run "git-town repo"
    Then "open" launches a new proposal with this url in my browser:
      """
      https://git.example.com/a/b
      """
