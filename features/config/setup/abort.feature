@messyoutput
Feature: aborting the setup assistant

  Background:
    Given a Git repo with origin
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG             | KEYS  |
      | welcome            | enter |
      | aliases            | enter |
      | main branch        | enter |
      | perennial branches | enter |
      | perennial regex    | esc   |

  Scenario: result
    Then Git Town runs no commands
    And the main branch is still not set
    And there are still no perennial branches
    And local Git setting "git-town.hosting-platform" still doesn't exist
    And local Git setting "git-town.push-new-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
