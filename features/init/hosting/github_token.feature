@messyoutput
Feature: enter the GitHub API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected GitHub platform
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG                | KEYS              | DESCRIPTION                                 |
      | welcome               | enter             |                                             |
      | aliases               | enter             |                                             |
      | main branch           | enter             |                                             |
      | perennial branches    |                   | no input here since the dialog doesn't show |
      | origin hostname       | enter             |                                             |
      | forge type            | enter             |                                             |
      | github connector type | enter             |                                             |
      | github token          | g h - t o k enter |                                             |
      | token scope           | enter             |                                             |
      | enter all             | enter             |                                             |
      | config storage        | enter             |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config git-town.github-token gh-tok  |
      | git config git-town.github-connector api |
    And local Git setting "git-town.github-token" is now "gh-tok"
    And local Git setting "git-town.forge-type" still doesn't exist

  Scenario: manually selected GitHub
    When I run "git-town init" and enter into the dialog:
      | DIALOG                | KEYS                                | DESCRIPTION                                 |
      | welcome               | enter                               |                                             |
      | aliases               | enter                               |                                             |
      | main branch           | enter                               |                                             |
      | perennial branches    |                                     | no input here since the dialog doesn't show |
      | origin hostname       | enter                               |                                             |
      | forge type            | down down down down down down enter |                                             |
      | github connector type | enter                               |                                             |
      | github token          | g h - t o k enter                   |                                             |
      | token scope           | enter                               |                                             |
      | enter all             | enter                               |                                             |
      | config storage        | enter                               |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config git-town.github-token gh-tok  |
      | git config git-town.forge-type github    |
      | git config git-town.github-connector api |
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.github-token" is now "gh-tok"

  Scenario: remove existing GitHub token
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And local Git setting "git-town.github-token" is "123"
    When I run "git-town init" and enter into the dialog:
      | DIALOG                | KEYS                                | DESCRIPTION                                 |
      | welcome               | enter                               |                                             |
      | aliases               | enter                               |                                             |
      | main branch           | enter                               |                                             |
      | perennial branches    |                                     | no input here since the dialog doesn't show |
      | origin hostname       | enter                               |                                             |
      | forge type            | enter                               |                                             |
      | github connector type | enter                               |                                             |
      | github token          | backspace backspace backspace enter |                                             |
      | enter all             | enter                               |                                             |
      | config storage        | enter                               |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config --unset git-town.github-token |
      | git config git-town.github-connector api |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist

  Scenario: store GitHub token globally
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    When I run "git-town init" and enter into the dialog:
      | DIALOG                | KEYS            | DESCRIPTION                                 |
      | welcome               | enter           |                                             |
      | aliases               | enter           |                                             |
      | main branch           | enter           |                                             |
      | perennial branches    |                 | no input here since the dialog doesn't show |
      | origin hostname       | enter           |                                             |
      | forge type            | enter           |                                             |
      | github connector type | enter           |                                             |
      | github token          | g h t o k enter |                                             |
      | token scope           | down enter      |                                             |
      | enter all             | enter           |                                             |
      | config storage        | enter           |                                             |
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config --global git-town.github-token ghtok |
      | git config git-town.github-connector api        |
    And global Git setting "git-town.github-token" is now "ghtok"

  Scenario: edit global GitHub token
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And global Git setting "git-town.github-token" is "123"
    When I run "git-town init" and enter into the dialog:
      | DIALOG                | KEYS                                      | DESCRIPTION                                 |
      | welcome               | enter                                     |                                             |
      | aliases               | enter                                     |                                             |
      | main branch           | enter                                     |                                             |
      | perennial branches    |                                           | no input here since the dialog doesn't show |
      | origin hostname       | enter                                     |                                             |
      | forge type            | enter                                     |                                             |
      | github connector type | enter                                     |                                             |
      | github token          | backspace backspace backspace 4 5 6 enter |                                             |
      | token scope           | enter                                     |                                             |
      | enter all             | enter                                     |                                             |
      | config storage        | enter                                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                       |
      | git config --global git-town.github-token 456 |
      | git config git-town.github-connector api      |
    And global Git setting "git-town.github-token" is now "456"
