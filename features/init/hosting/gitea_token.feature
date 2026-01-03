@messyoutput
Feature: enter the Gitea API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected Gitea platform
    And my repo's "origin" remote is "git@gitea.com:git-town/git-town.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                    | DESCRIPTION                                 |
      | welcome            | enter                   |                                             |
      | aliases            | enter                   |                                             |
      | main branch        | enter                   |                                             |
      | perennial branches |                         | no input here since the dialog doesn't show |
      | origin hostname    | enter                   |                                             |
      | forge type         | enter                   | auto-detect                                 |
      | gitea token        | g i t e a - t o k enter |                                             |
      | token scope        | enter                   |                                             |
      | enter all          | enter                   |                                             |
      | config storage     | enter                   | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                   |
      | git config git-town.gitea-token gitea-tok |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.gitea-token" is now "gitea-tok"

  Scenario: select Gitea manually
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                           | DESCRIPTION                                 |
      | welcome            | enter                          |                                             |
      | aliases            | enter                          |                                             |
      | main branch        | enter                          |                                             |
      | perennial branches |                                | no input here since the dialog doesn't show |
      | origin hostname    | enter                          |                                             |
      | forge type         | down down down down down enter |                                             |
      | gitea token        | g i t e a - t o k enter        |                                             |
      | token scope        | enter                          |                                             |
      | enter all          | enter                          |                                             |
      | config storage     | enter                          | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                   |
      | git config git-town.gitea-token gitea-tok |
      | git config git-town.forge-type gitea      |
    And local Git setting "git-town.forge-type" is now "gitea"
    And local Git setting "git-town.gitea-token" is now "gitea-tok"

  Scenario: store Gitea API token globally
    And my repo's "origin" remote is "git@gitea.com:git-town/git-town.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                    | DESCRIPTION                                 |
      | welcome            | enter                   |                                             |
      | aliases            | enter                   |                                             |
      | main branch        | enter                   |                                             |
      | perennial branches |                         | no input here since the dialog doesn't show |
      | origin hostname    | enter                   |                                             |
      | forge type         | enter                   |                                             |
      | gitea token        | g i t e a - t o k enter |                                             |
      | token scope        | down enter              |                                             |
      | enter all          | enter                   |                                             |
      | config storage     | enter                   | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                            |
      | git config --global git-town.gitea-token gitea-tok |
    And global Git setting "git-town.gitea-token" is now "gitea-tok"

  Scenario: edit global Gitea token
    Given my repo's "origin" remote is "git@gitea.com:git-town/git-town.git"
    And global Git setting "git-town.gitea-token" is "123"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                                      | DESCRIPTION                                 |
      | welcome            | enter                                     |                                             |
      | aliases            | enter                                     |                                             |
      | main branch        | enter                                     |                                             |
      | perennial branches |                                           | no input here since the dialog doesn't show |
      | origin hostname    | enter                                     |                                             |
      | forge type         | enter                                     |                                             |
      | gitea token        | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope        | enter                                     |                                             |
      | enter all          | enter                                     |                                             |
      | config storage     | enter                                     | git metadata                                |
    Then Git Town runs the commands
      | COMMAND                                      |
      | git config --global git-town.gitea-token 456 |
    And global Git setting "git-town.gitea-token" is now "456"
