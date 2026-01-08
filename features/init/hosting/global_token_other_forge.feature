@messyoutput
Feature: a global API token of another forge exists

  Scenario: on GitHub, with global GitLab token
    Given a Git repo with origin
    And my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And global Git setting "git-town.gitlab-token" is "987654"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS            | DESCRIPTION                                 |
      | welcome            | enter           |                                             |
      | aliases            | enter           |                                             |
      | main branch        | enter           |                                             |
      | perennial branches |                 | no input here since the dialog doesn't show |
      | origin hostname    | enter           |                                             |
      | forge type         | enter           |                                             |
      | github connector   | enter           |                                             |
      | github token       | g h t o k enter |                                             |
      | token scope        | enter           |                                             |
      | enter all          | enter           |                                             |
      | config storage     | enter           |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config git-town.github-token ghtok   |
      | git config git-town.github-connector api |
    And global Git setting "git-town.gitlab-token" is still "987654"
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" is now "ghtok"
    And local Git setting "git-town.gitlab-token" now doesn't exist
