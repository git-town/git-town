@messyoutput
Feature: enter the GitHub API token

  Background:
    Given a Git repo with origin

  Scenario: auto-detected GitHub platform
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main branch                 | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | perennial regex             | enter             |                                             |
      | default branch type         | enter             |                                             |
      | feature regex               | enter             |                                             |
      | dev-remote                  | enter             |                                             |
      | forge type: auto-detect     | enter             |                                             |
      | github token                | 1 2 3 4 5 6 enter |                                             |
      | origin hostname             | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-prototype-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | sync-tags                   | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | new-branch-type             | down enter        |                                             |
      | ship-strategy               | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | save config to Git metadata | down enter        |                                             |
    Then Git Town runs the commands
      | COMMAND                                 |
      | git config git-town.github-token 123456 |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" is now "123456"

  Scenario: manually selected GitHub
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                           | DESCRIPTION                                 |
      | welcome                     | enter                          |                                             |
      | aliases                     | enter                          |                                             |
      | main branch                 | enter                          |                                             |
      | perennial branches          |                                | no input here since the dialog doesn't show |
      | perennial regex             | enter                          |                                             |
      | default branch type         | enter                          |                                             |
      | feature regex               | enter                          |                                             |
      | dev-remote                  | enter                          |                                             |
      | forge type                  | down down down down down enter |                                             |
      | github token                | 1 2 3 4 5 6 enter              |                                             |
      | origin hostname             | enter                          |                                             |
      | sync-feature-strategy       | enter                          |                                             |
      | sync-perennial-strategy     | enter                          |                                             |
      | sync-prototype-strategy     | enter                          |                                             |
      | sync-upstream               | enter                          |                                             |
      | sync-tags                   | enter                          |                                             |
      | push-new-branches           | enter                          |                                             |
      | push-hook                   | enter                          |                                             |
      | new-branch-type             | enter                          |                                             |
      | ship-strategy               | enter                          |                                             |
      | ship-delete-tracking-branch | enter                          |                                             |
      | save config to Git metadata | down enter                     |                                             |
    Then Git Town runs the commands
      | COMMAND                                 |
      | git config git-town.github-token 123456 |
      | git config git-town.forge-type github   |
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.github-token" is now "123456"

  Scenario: remove existing GitHub token
    Given my repo's "origin" remote is "git@github.com:git-town/git-town.git"
    And local Git setting "git-town.github-token" is "123"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS                                | DESCRIPTION                                 |
      | welcome                     | enter                               |                                             |
      | aliases                     | enter                               |                                             |
      | main branch                 | enter                               |                                             |
      | perennial branches          |                                     | no input here since the dialog doesn't show |
      | perennial regex             | enter                               |                                             |
      | default branch type         | enter                               |                                             |
      | feature regex               | enter                               |                                             |
      | dev-remote                  | enter                               |                                             |
      | forge type: auto-detect     | enter                               |                                             |
      | github token                | backspace backspace backspace enter |                                             |
      | origin hostname             | enter                               |                                             |
      | sync-feature-strategy       | enter                               |                                             |
      | sync-perennial-strategy     | enter                               |                                             |
      | sync-prototype-strategy     | enter                               |                                             |
      | sync-upstream               | enter                               |                                             |
      | sync-tags                   | enter                               |                                             |
      | push-new-branches           | enter                               |                                             |
      | push-hook                   | enter                               |                                             |
      | new-branch-type             | enter                               |                                             |
      | ship-strategy               | enter                               |                                             |
      | ship-delete-tracking-branch | enter                               |                                             |
      | save config to Git metadata | down enter                          |                                             |
    Then Git Town runs the commands
      | COMMAND                                  |
      | git config --unset git-town.github-token |
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist

  Scenario: undo
    When I run "git-town undo"
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist
