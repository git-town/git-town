Feature: enter Git Town configuration

  Scenario: unconfigured, accept all default values --> working setup
    Given the branches "dev" and "production"
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS  |
      | aliases                     | enter |
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
      | add all aliases                           | a enter                |
      | accept the already configured main branch | enter                  |
      | configure the perennial branches          | space down space enter |
      | sync-feature-strategy                     | down enter             |
      | sync-perennial-strategy                   | down enter             |
      | sync-upstream                             | down enter             |
      | enable push-new-branches                  | down enter             |
      | disable the push hook                     | down enter             |
      | disable ship-delete-tracking-branch       | down enter             |
      | sync-before-ship                          | down enter             |
    Then global Git setting "alias.append" is now "town append"
    And global Git setting "alias.diff-parent" is now "town diff-parent"
    And global Git setting "alias.hack" is now "town hack"
    And global Git setting "alias.kill" is now "town kill"
    And global Git setting "alias.prepend" is now "town prepend"
    And global Git setting "alias.propose" is now "town propose"
    And global Git setting "alias.rename-branch" is now "town rename-branch"
    And global Git setting "alias.repo" is now "town repo"
    And global Git setting "alias.set-parent" is now "town set-parent"
    And global Git setting "alias.ship" is now "town ship"
    And global Git setting "alias.sync" is now "town sync"
    And the main branch is now "main"
    And the perennial branches are now "production"
    And local Git Town setting "sync-feature-strategy" is now "rebase"
    And local Git Town setting "sync-perennial-strategy" is now "merge"
    And local Git Town setting "sync-upstream" is now "false"
    And local Git Town setting "push-new-branches" is now "true"
    And local Git Town setting "push-hook" is now "true"
    And local Git Town setting "ship-delete-tracking-branch" is now "false"
    And local Git Town setting "sync-before-ship" is now "true"

  Scenario: override an existing Git alias
    Given I ran "git config --global alias.append checkout"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS    |
      | aliases                     | o enter |
      | main development branch     | enter   |
      | perennial branches          | enter   |
      | sync-feature-strategy       | enter   |
      | sync-perennial-strategy     | enter   |
      | sync-upstream               | enter   |
      | push-new-branches           | enter   |
      | push-hook                   | enter   |
      | ship-delete-tracking-branch | enter   |
      | sync-before-ship            | enter   |
    Then global Git setting "alias.append" is now "town append"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given Git Town is not configured
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS       | DESCRIPTION                                 |
      | aliases                     | enter      |                                             |
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
