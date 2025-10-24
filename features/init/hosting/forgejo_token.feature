@messyoutput
Feature: enter the Forgejo API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected Forgejo platform
    Given my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                   | DESCRIPTION                                 |
      | welcome            | enter                  |                                             |
      | aliases            | enter                  |                                             |
      | main branch        | enter                  |                                             |
      | perennial branches |                        | no input here since the dialog doesn't show |
      | origin hostname    | enter                  |                                             |
      | forge type         | enter                  |                                             |
      | forgejo token      | c o d e - t o k  enter |                                             |
      | token scope        | enter                  |                                             |
      | enter all          | enter                  |                                             |
      | config storage     | enter                  |                                             |
    Then Git Town runs the commands
      | COMMAND                                    |
      | git config git-town.forgejo-token code-tok |
    And local Git setting "git-town.forgejo-token" is now "code-tok"
    And local Git setting "git-town.forge-type" still doesn't exist

  Scenario: select Forgejo manually
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                      | DESCRIPTION                                 |
      | welcome            | enter                     |                                             |
      | aliases            | enter                     |                                             |
      | main branch        | enter                     |                                             |
      | perennial branches |                           | no input here since the dialog doesn't show |
      | origin hostname    | enter                     |                                             |
      | forge type         | down down down down enter |                                             |
      | forgejo token      | c o d e - t o k  enter    |                                             |
      | token scope        | enter                     |                                             |
      | enter all          | enter                     |                                             |
      | config storage     | enter                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                    |
      | git config git-town.forgejo-token code-tok |
      | git config git-town.forge-type forgejo     |
    And local Git setting "git-town.forge-type" is now "forgejo"
    And local Git setting "git-town.forgejo-token" is now "code-tok"

  Scenario: store Forgejo API token globally
    And my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                   | DESCRIPTION                                 |
      | welcome            | enter                  |                                             |
      | aliases            | enter                  |                                             |
      | main-branch        | enter                  |                                             |
      | perennial branches |                        | no input here since the dialog doesn't show |
      | origin hostname    | enter                  |                                             |
      | forge type         | enter                  |                                             |
      | forgejo token      | c o d e - t o k  enter |                                             |
      | token scope        | down enter             |                                             |
      | enter all          | enter                  |                                             |
      | config storage     | enter                  |                                             |
    Then Git Town runs the commands
      | COMMAND                                             |
      | git config --global git-town.forgejo-token code-tok |
    And global Git setting "git-town.forgejo-token" is now "code-tok"

  Scenario: edit global Forgejo API token
    And my repo's "origin" remote is "git@codeberg.org:git-town/docs.git"
    Given global Git setting "git-town.forgejo-token" is "code123"
    When I run "git-town init" and enter into the dialog:
      | DIALOG             | KEYS                                      | DESCRIPTION                                 |
      | welcome            | enter                                     |                                             |
      | aliases            | enter                                     |                                             |
      | main branch        | enter                                     |                                             |
      | perennial branches |                                           | no input here since the dialog doesn't show |
      | origin hostname    | enter                                     |                                             |
      | forge type         | enter                                     |                                             |
      | forgejo token      | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope        | enter                                     |                                             |
      | enter all          | enter                                     |                                             |
      | config storage     | enter                                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                            |
      | git config --global git-town.forgejo-token code456 |
    And global Git setting "git-town.forgejo-token" is now "code456"
