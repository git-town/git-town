@messyoutput
Feature: enter the GitLab API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected GitLab platform
    Given my repo's "origin" remote is "git@gitlab.com:git-town/git-town.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS              | DESCRIPTION                                 |
      | welcome            | enter             |                                             |
      | aliases            | enter             |                                             |
      | main branch        | enter             |                                             |
      | perennial branches |                   | no input here since the dialog doesn't show |
      | origin hostname    | enter             |                                             |
      | forge type         | enter             |                                             |
      | gitlab connector   | enter             |                                             |
      | gitlab token       | g l - t o k enter |                                             |
      | token scope        | enter             |                                             |
      | enter all          | enter             |                                             |
      | config storage     | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config git-town.gitlab-token gl-tok  |
      | git config git-town.gitlab-connector api |
    And local Git setting "git-town.forge-type" still doesn't exist

  Scenario: select GitLab manually
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS              | DESCRIPTION                                 |
      | welcome            | enter             |                                             |
      | aliases            | enter             |                                             |
      | main branch        | enter             |                                             |
      | perennial branches |                   | no input here since the dialog doesn't show |
      | origin hostname    | enter             |                                             |
      | forge type         | up enter          |                                             |
      | gitlab connector   | enter             |                                             |
      | gitlab token       | g l - t o k enter |                                             |
      | token scope        | enter             |                                             |
      | enter all          | enter             |                                             |
      | config storage     | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config git-town.gitlab-token gl-tok  |
      | git config git-town.forge-type gitlab    |
      | git config git-town.gitlab-connector api |
    And local Git setting "git-town.forge-type" is now "gitlab"
    And local Git setting "git-town.gitlab-token" is now "gl-tok"

  Scenario: store GitLab API token globally
    Given my repo's "origin" remote is "git@gitlab.com:git-town/git-town.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS            | DESCRIPTION                                 |
      | welcome            | enter           |                                             |
      | aliases            | enter           |                                             |
      | main branch        | enter           |                                             |
      | perennial branches |                 | no input here since the dialog doesn't show |
      | origin hostname    | enter           |                                             |
      | forge type         | enter           |                                             |
      | gitlab connector   | enter           | api                                         |
      | gitlab token       | g l t o k enter |                                             |
      | token scope        | down enter      |                                             |
      | enter all          | enter           |                                             |
      | config storage     | enter           | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --global git-town.gitlab-token gltok |
      | git config git-town.gitlab-connector api        |
    And global Git setting "git-town.gitlab-token" is now "gltok"

  Scenario: edit global GitLab token
    Given my repo's "origin" remote is "git@gitlab.com:git-town/git-town.git"
    And global Git setting "git-town.gitlab-token" is "123"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                                      | DESCRIPTION                                 |
      | welcome            | enter                                     |                                             |
      | aliases            | enter                                     |                                             |
      | main branch        | enter                                     |                                             |
      | perennial branches |                                           | no input here since the dialog doesn't show |
      | origin hostname    | enter                                     |                                             |
      | forge type         | enter                                     |                                             |
      | gitlab connector   | enter                                     |                                             |
      | gitlab token       | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope        | enter                                     |                                             |
      | enter all          | enter                                     |                                             |
      | config storage     | enter                                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                       |
      | git config --global git-town.gitlab-token 456 |
      | git config git-town.gitlab-connector api      |
    And global Git setting "git-town.gitlab-token" is now "456"
