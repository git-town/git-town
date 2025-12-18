@messyoutput
Feature: aborting the setup assistant

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And local Git setting "init.defaultbranch" is "main"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG          | KEYS  |
      | welcome         | enter |
      | aliases         | enter |
      | main branch     | enter |
      | origin hostname | esc   |

  Scenario: result
    Then Git Town runs no commands
    # keep-sorted start
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
    And the main branch is still not set
    And there are still no perennial branches
