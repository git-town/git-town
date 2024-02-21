Feature: change existing information in Git metadata

  Background:
    Given a perennial branch "qa"
    And a branch "production"
    And the main branch is "main"
    And local Git Town setting "push-new-branches" is "false"
    And local Git Town setting "push-hook" is "false"
    When I run "git-town config setup" and enter into the dialogs:
      | DESCRIPTION                               | KEYS                   |
      | welcome                                   | enter                  |
      | add all aliases                           | a enter                |
      | accept the already configured main branch | enter                  |
      | change the perennial branches             | space down space enter |
      | enter a perennial regex                   | 3 3 6 6 enter          |
      | set github as hosting service             | up up enter            |
      | github token                              | 1 2 3 4 5 6 enter      |
      | origin hostname                           | c o d e enter          |
      | sync-feature-strategy                     | down enter             |
      | sync-perennial-strategy                   | down enter             |
      | sync-upstream                             | down enter             |
      | enable push-new-branches                  | down enter             |
      | disable the push hook                     | down enter             |
      | disable ship-delete-tracking-branch       | down enter             |
      | sync-before-ship                          | down enter             |
      | save config to Git metadata               | down enter             |

  Scenario: result
    Then it runs the commands
      | COMMAND                                                      |
      | git config --global alias.append "town append"               |
      | git config --global alias.diff-parent "town diff-parent"     |
      | git config --global alias.hack "town hack"                   |
      | git config --global alias.kill "town kill"                   |
      | git config --global alias.prepend "town prepend"             |
      | git config --global alias.propose "town propose"             |
      | git config --global alias.rename-branch "town rename-branch" |
      | git config --global alias.repo "town repo"                   |
      | git config --global alias.set-parent "town set-parent"       |
      | git config --global alias.ship "town ship"                   |
      | git config --global alias.sync "town sync"                   |
      | git config git-town.github-token 123456                      |
      | git config git-town.hosting-platform github                  |
      | git config git-town.hosting-origin-hostname code             |
    And global Git setting "alias.append" is now "town append"
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
    And local Git Town setting "hosting-platform" is now "github"
    And local Git Town setting "github-token" is now "123456"
    And local Git Town setting "hosting-origin-hostname" is now "code"
    And local Git Town setting "sync-feature-strategy" is now "rebase"
    And local Git Town setting "sync-perennial-strategy" is now "merge"
    And local Git Town setting "sync-upstream" is now "false"
    And local Git Town setting "perennial-regex" is now "3366"
    And local Git Town setting "push-new-branches" is now "true"
    And local Git Town setting "push-hook" is now "true"
    And local Git Town setting "ship-delete-tracking-branch" is now "false"
    And local Git Town setting "sync-before-ship" is now "true"
