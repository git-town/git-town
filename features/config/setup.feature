Feature: enter Git Town configuration

  Scenario: unconfigured, accept all default values --> working setup
    Given the branches "dev" and "production"
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS  |
      | main development branch     | enter |
      | perennial branches          | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-upstream               | enter |
      | push-new-branches           | enter |
      | push-hook                   | enter |
      | ship-delete-tracking-branch | enter |
      | sync-before-ship            | enter |
    Then the main branch is now "main"
    And there are still no perennial branches
    And local Git Town setting "push-new-branches" is now "false"
    And local Git Town setting "push-hook" is now "true"
    And local Git Town setting "sync-feature-strategy" is now "merge"
    And local Git Town setting "sync-perennial-strategy" is now "rebase"
    And local Git Town setting "sync-upstream" is now "true"
    And local Git Town setting "ship-delete-tracking-branch" is now "true"
    And local Git Town setting "sync-before-ship" is now "false"

  Scenario: change existing configuration
    Given a perennial branch "qa"
    And a branch "production"
    And the main branch is "main"
    And local Git Town setting "push-new-branches" is "false"
    And local Git Town setting "push-hook" is "false"
    When I run "git-town config setup" and enter into the dialogs:
      | DESCRIPTION                               | KEYS                   |
      | accept the already configured main branch | enter                  |
      | configure the perennial branches          | space down space enter |
      | sync-feature-strategy                     | down enter             |
      | sync-perennial-strategy                   | down enter             |
      | sync-upstream                             | down enter             |
      | enable push-new-branches                  | down enter             |
      | disable the push hook                     | down enter             |
      | disable ship-delete-tracking-branch       | down enter             |
      | sync-before-ship                          | down enter             |
    Then the main branch is now "main"
    And the perennial branches are now "production"
    And local Git Town setting "sync-feature-strategy" is now "rebase"
    And local Git Town setting "sync-perennial-strategy" is now "merge"
    And local Git Town setting "sync-upstream" is now "false"
    And local Git Town setting "push-new-branches" is now "true"
    And local Git Town setting "push-hook" is now "true"
    And local Git Town setting "ship-delete-tracking-branch" is now "false"
    And local Git Town setting "sync-before-ship" is now "true"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given Git Town is not configured
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS       | DESCRIPTION                                 |
      | main development branch     | down enter |                                             |
      | perennial branches          |            | no input here since the dialog doesn't show |
      | sync-feature-strategy       | enter      |                                             |
      | sync-perennial-strategy     | enter      |                                             |
      | sync-upstream               | enter      |                                             |
      | push-new-branches           | enter      |                                             |
      | push-hook                   | enter      |                                             |
      | ship-delete-tracking-branch | enter      |                                             |
      | sync-before-ship            | enter      |                                             |
    Then the main branch is now "main"
    And there are still no perennial branches
