Feature: remove existing configuration in Git metadata

  Background:
    Given a perennial branch "qa"
    And a branch "production"
    And the main branch is "main"
    And global Git setting "alias.append" is "town append"
    And global Git setting "alias.diff-parent" is "town diff-parent"
    And global Git setting "alias.hack" is "town hack"
    And global Git setting "alias.kill" is "town kill"
    And global Git setting "alias.prepend" is "town prepend"
    And global Git setting "alias.propose" is "town propose"
    And global Git setting "alias.rename-branch" is "town rename-branch"
    And global Git setting "alias.repo" is "town repo"
    And global Git setting "alias.set-parent" is "town set-parent"
    And global Git setting "alias.ship" is "town ship"
    And global Git setting "alias.sync" is "town sync"
    And local Git Town setting "hosting-platform" is "github"
    And local Git Town setting "push-new-branches" is "false"
    And local Git Town setting "push-hook" is "false"
    And local Git Town setting "hosting-origin-hostname" is "code"
    And local Git Town setting "sync-feature-strategy" is "rebase"
    And local Git Town setting "sync-perennial-strategy" is "rebase"
    And local Git Town setting "sync-upstream" is "true"
    And local Git Town setting "push-new-branches" is "true"
    And local Git Town setting "push-hook" is "true"
    And local Git Town setting "ship-delete-tracking-branch" is "false"
    And local Git Town setting "sync-before-ship" is "true"
    When I run "git-town config setup" and enter into the dialogs:
      | DESCRIPTION                             | KEYS                                          |
      | welcome                                 | enter                                         |
      | add all aliases                         | n enter                                       |
      | keep the already configured main branch | enter                                         |
      | change the perennial branches           | space down space enter                        |
      | remove hosting service override         | up up up enter                                |
      | remove origin hostname                  | backspace backspace backspace backspace enter |
      | sync-feature-strategy                   | down enter                                    |
      | sync-perennial-strategy                 | down enter                                    |
      | sync-upstream                           | down enter                                    |
      | enable push-new-branches                | down enter                                    |
      | disable the push hook                   | down enter                                    |
      | disable ship-delete-tracking-branch     | down enter                                    |
      | sync-before-ship                        | down enter                                    |
      | save config to Git metadata             | down enter                                    |

  Scenario: result
    Then it runs the commands
      | COMMAND                                             |
      | git config --global --unset alias.append            |
      | git config --global --unset alias.diff-parent       |
      | git config --global --unset alias.hack              |
      | git config --global --unset alias.kill              |
      | git config --global --unset alias.prepend           |
      | git config --global --unset alias.propose           |
      | git config --global --unset alias.rename-branch     |
      | git config --global --unset alias.repo              |
      | git config --global --unset alias.set-parent        |
      | git config --global --unset alias.ship              |
      | git config --global --unset alias.sync              |
      | git config --unset git-town.hosting-platform        |
      | git config --unset git-town.hosting-origin-hostname |
    And global Git setting "alias.append" now doesn't exist
    And global Git setting "alias.diff-parent" now doesn't exist
    And global Git setting "alias.hack" now doesn't exist
    And global Git setting "alias.kill" now doesn't exist
    And global Git setting "alias.prepend" now doesn't exist
    And global Git setting "alias.propose" now doesn't exist
    And global Git setting "alias.rename-branch" now doesn't exist
    And global Git setting "alias.repo" now doesn't exist
    And global Git setting "alias.set-parent" now doesn't exist
    And global Git setting "alias.ship" now doesn't exist
    And global Git setting "alias.sync" now doesn't exist
    And the main branch is still "main"
    And the perennial branches are now "production"
    And local Git Town setting "hosting-platform" now doesn't exist
    And local Git Town setting "github-token" now doesn't exist
    And local Git Town setting "hosting-origin-hostname" now doesn't exist
    And local Git Town setting "sync-feature-strategy" is now "merge"
    And local Git Town setting "sync-perennial-strategy" is now "merge"
    And local Git Town setting "sync-upstream" is now "false"
    And local Git Town setting "push-new-branches" is now "false"
    And local Git Town setting "push-hook" is now "false"
    And local Git Town setting "ship-delete-tracking-branch" is now "true"
    And local Git Town setting "sync-before-ship" is now "false"
