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
      | hosting platform            | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-upstream               | enter |
      | push-new-branches           | enter |
      | push-hook                   | enter |
      | ship-delete-tracking-branch | enter |
      | sync-before-ship            | enter |
    Then it runs no commands
    And the main branch is now "main"
    And there are still no perennial branches
    And local Git Town setting "code-hosting-platform" is still not set
    And local Git Town setting "push-new-branches" is still not set
    And local Git Town setting "push-hook" is still not set
    And local Git Town setting "sync-feature-strategy" is still not set
    And local Git Town setting "sync-perennial-strategy" is still not set
    And local Git Town setting "sync-upstream" is still not set
    And local Git Town setting "ship-delete-tracking-branch" is still not set
    And local Git Town setting "sync-before-ship" is still not set

  Scenario: unconfigured, enter some values and hit ESC --> does not save
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                  | KEYS  |
      | aliases                 | enter |
      | main development branch | enter |
      | perennial branches      | enter |
      | hosting platform        | esc   |
    Then it runs no commands
    And the main branch is still not set
    And there are still no perennial branches
    And local Git Town setting "code-hosting-platform" is still not set
    And local Git Town setting "push-new-branches" is still not set
    And local Git Town setting "push-hook" is still not set
    And local Git Town setting "sync-feature-strategy" is still not set
    And local Git Town setting "sync-perennial-strategy" is still not set
    And local Git Town setting "sync-upstream" is still not set
    And local Git Town setting "ship-delete-tracking-branch" is still not set
    And local Git Town setting "sync-before-ship" is still not set

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
      | set github as hosting service             | up up enter            |
      | github token                              | 1 2 3 4 5 6 enter      |
      | sync-feature-strategy                     | down enter             |
      | sync-perennial-strategy                   | down enter             |
      | sync-upstream                             | down enter             |
      | enable push-new-branches                  | down enter             |
      | disable the push hook                     | down enter             |
      | disable ship-delete-tracking-branch       | down enter             |
      | sync-before-ship                          | down enter             |
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
      | git config git-town.code-hosting-platform github             |
      | git config git-town.github-token 123456                      |
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
    And local Git Town setting "code-hosting-platform" is now "github"
    And local Git Town setting "github-token" is now "123456"
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
      | hosting platform            | enter   |
      | sync-feature-strategy       | enter   |
      | sync-perennial-strategy     | enter   |
      | sync-upstream               | enter   |
      | push-new-branches           | enter   |
      | push-hook                   | enter   |
      | ship-delete-tracking-branch | enter   |
      | sync-before-ship            | enter   |
    Then it runs the commands
      | COMMAND                                        |
      | git config --global alias.append "town append" |
    And global Git setting "alias.append" is now "town append"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given Git Town is not configured
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS       | DESCRIPTION                                 |
      | aliases                     | enter      |                                             |
      | main development branch     | down enter |                                             |
      | perennial branches          |            | no input here since the dialog doesn't show |
      | hosting platform            | enter      |                                             |
      | sync-feature-strategy       | enter      |                                             |
      | sync-perennial-strategy     | enter      |                                             |
      | sync-upstream               | enter      |                                             |
      | push-new-branches           | enter      |                                             |
      | push-hook                   | enter      |                                             |
      | ship-delete-tracking-branch | enter      |                                             |
      | sync-before-ship            | enter      |                                             |
    Then the main branch is now "main"
    And there are still no perennial branches

  Scenario: remove an existing code hosting override
    Given local Git Town setting "code-hosting-platform" is "github"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS           | DESCRIPTION                                 |
      | aliases                     | enter          |                                             |
      | main development branch     | down enter     |                                             |
      | perennial branches          |                | no input here since the dialog doesn't show |
      | hosting platform            | up up up enter |                                             |
      | sync-feature-strategy       | enter          |                                             |
      | sync-perennial-strategy     | enter          |                                             |
      | sync-upstream               | enter          |                                             |
      | push-new-branches           | enter          |                                             |
      | push-hook                   | enter          |                                             |
      | ship-delete-tracking-branch | enter          |                                             |
      | sync-before-ship            | enter          |                                             |
    Then it runs the commands
      | COMMAND                                           |
      | git config --unset git-town.code-hosting-platform |
    And local Git Town setting "code-hosting-platform" is now not set

  Scenario: enter a GitLab token
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | aliases                     | enter             |                                             |
      | main development branch     | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | hosting platform            | up enter          |                                             |
      | gitlab token                | 1 2 3 4 5 6 enter |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | sync-before-ship            | enter             |                                             |
    Then it runs the commands
      | COMMAND                                          |
      | git config git-town.code-hosting-platform gitlab |
      | git config git-town.gitlab-token 123456          |
    And local Git Town setting "code-hosting-platform" is now "gitlab"
    And local Git Town setting "gitlab-token" is now "123456"

  Scenario: enter a Gitea token
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | aliases                     | enter             |                                             |
      | main development branch     | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | hosting platform            | down down enter   |                                             |
      | gitea token                 | 1 2 3 4 5 6 enter |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | sync-before-ship            | enter             |                                             |
    Then it runs the commands
      | COMMAND                                         |
      | git config git-town.code-hosting-platform gitea |
      | git config git-town.gitea-token 123456          |
    And local Git Town setting "code-hosting-platform" is now "gitea"
    And local Git Town setting "gitea-token" is now "123456"
