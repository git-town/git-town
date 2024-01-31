Feature: enter Git Town configuration

  Scenario: unconfigured, accept all default values --> working setup with config file
    Given the branches "dev" and "production"
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS  |
      | welcome                     | enter |
      | aliases                     | enter |
      | main development branch     | enter |
      | perennial branches          | enter |
      | hosting platform            | enter |
      | origin hostname             | enter |
      | sync-feature-strategy       | enter |
      | sync-perennial-strategy     | enter |
      | sync-upstream               | enter |
      | push-new-branches           | enter |
      | push-hook                   | enter |
      | ship-delete-tracking-branch | enter |
      | sync-before-ship            | enter |
      | save config file            | enter |
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
    And the configuration file is now:
      """
      # Git Town configuration file
      #
      # The "push-hook" setting determines whether Git Town
      # permits or prevents Git hooks while pushing branches.
      # Hooks are enabled by default. If your Git hooks are slow,
      # you can disable them to speed up branch syncing.
      #
      # When disabled, Git Town pushes using the "--no-verify" switch.
      # More info at https://www.git-town.com/preferences/push-hook.
      push-hook = true

      # Should Git Town push the new branches it creates
      # immediately to origin even if they are empty?
      #
      # When enabled, you can run "git push" right away
      # but creating new branches is slower and
      # it triggers an unnecessary CI run on the empty branch.
      #
      # When disabled, many Git Town commands execute faster
      # and Git Town will create the missing tracking branch
      # on the first run of "git sync".
      push-new-branches = false

      # Should "git ship" delete the tracking branch?
      # You want to disable this if your code hosting system
      # (GitHub, GitLab, etc) deletes head branches when
      # merging pull requests through its UI.
      ship-delete-tracking-branch = true

      # Should "git ship" sync branches before shipping them?
      #
      # Guidance: enable when shipping branches locally on your machine
      # and disable when shipping feature branches via the code hosting
      # API or web UI.
      #
      # When enabled, branches are always fully up to date when shipped
      # and you get a chance to resolve merge conflicts
      # between the feature branch to ship and the main development branch
      # on the feature branch. This helps keep the main branch green.
      # But this also triggers another CI run and delays shipping.
      sync-before-ship = false

      # Should "git sync" also fetch updates from the upstream remote?
      #
      # If an "upstream" remote exists, and this setting is enabled,
      # "git sync" will also update the local main branch
      # with commits from the main branch at the upstream remote.
      #
      # This is useful if the repository you work on is a fork,
      # and you want to keep it in sync with the repo it was forked from.
      sync-upstream = true

      [branches]

        # The main branch is the branch from which you cut new feature branches,
        # and into which you ship feature branches when they are done.
        # This branch is often called "main", "master", or "development".
        main = "main"

        # Perennial branches are long-lived branches.
        # They are never shipped and have no ancestors.
        # Typically, perennial branches have names like
        # "development", "staging", "qa", "production", etc.
        perennials = []

      [sync-strategy]

        # How should Git Town synchronize feature branches?
        # Feature branches are short-lived branches cut from
        # the main branch and shipped back into the main branch.
        # Typically you develop features and bug fixes on them,
        # hence their name.
        feature-branches = "merge"

        # How should Git Town synchronize perennial branches?
        # Perennial branches have no parent branch.
        # The only updates they receive are additional commits
        # made to their tracking branch somewhere else.
        perennial-branches = "rebase"
      """

  Scenario: unconfigured, enter some values and hit ESC --> does not save
    And local Git setting "init.defaultbranch" is "main"
    And Git Town is not configured
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                  | KEYS  |
      | welcome                 | enter |
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

  Scenario: change existing configuration in Git metadata
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
      | save config file                          | down enter             |
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
      | git config git-town.code-hosting-platform github             |
      | git config git-town.code-hosting-origin-hostname code        |
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
    And local Git Town setting "code-hosting-origin-hostname" is now "code"
    And local Git Town setting "sync-feature-strategy" is now "rebase"
    And local Git Town setting "sync-perennial-strategy" is now "merge"
    And local Git Town setting "sync-upstream" is now "false"
    And local Git Town setting "push-new-branches" is now "true"
    And local Git Town setting "push-hook" is now "true"
    And local Git Town setting "ship-delete-tracking-branch" is now "false"
    And local Git Town setting "sync-before-ship" is now "true"

  Scenario: remove existing configuration in Git metadata
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
    And local Git Town setting "code-hosting-platform" is "github"
    And local Git Town setting "push-new-branches" is "false"
    And local Git Town setting "push-hook" is "false"
    And local Git Town setting "code-hosting-origin-hostname" is "code"
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
      | save config file                        | down enter                                    |
    Then it runs the commands
      | COMMAND                                                  |
      | git config --global --unset alias.append                 |
      | git config --global --unset alias.diff-parent            |
      | git config --global --unset alias.hack                   |
      | git config --global --unset alias.kill                   |
      | git config --global --unset alias.prepend                |
      | git config --global --unset alias.propose                |
      | git config --global --unset alias.rename-branch          |
      | git config --global --unset alias.repo                   |
      | git config --global --unset alias.set-parent             |
      | git config --global --unset alias.ship                   |
      | git config --global --unset alias.sync                   |
      | git config --unset git-town.code-hosting-platform        |
      | git config --unset git-town.code-hosting-origin-hostname |
    And global Git setting "alias.append" no longer exists
    And global Git setting "alias.diff-parent" no longer exists
    And global Git setting "alias.hack" no longer exists
    And global Git setting "alias.kill" no longer exists
    And global Git setting "alias.prepend" no longer exists
    And global Git setting "alias.propose" no longer exists
    And global Git setting "alias.rename-branch" no longer exists
    And global Git setting "alias.repo" no longer exists
    And global Git setting "alias.set-parent" no longer exists
    And global Git setting "alias.ship" no longer exists
    And global Git setting "alias.sync" no longer exists
    And the main branch is still "main"
    And the perennial branches are now "production"
    And local Git Town setting "code-hosting-platform" no longer exists
    And local Git Town setting "github-token" no longer exists
    And local Git Town setting "code-hosting-origin-hostname" no longer exists
    And local Git Town setting "sync-feature-strategy" is now "merge"
    And local Git Town setting "sync-perennial-strategy" is now "merge"
    And local Git Town setting "sync-upstream" is now "false"
    And local Git Town setting "push-new-branches" is now "false"
    And local Git Town setting "push-hook" is now "false"
    And local Git Town setting "ship-delete-tracking-branch" is now "true"
    And local Git Town setting "sync-before-ship" is now "false"

  Scenario: override an existing Git alias
    Given I ran "git config --global alias.append checkout"
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS    |
      | welcome                     | enter   |
      | aliases                     | o enter |
      | main development branch     | enter   |
      | perennial branches          | enter   |
      | hosting platform            | enter   |
      | origin hostname             | enter   |
      | sync-feature-strategy       | enter   |
      | sync-perennial-strategy     | enter   |
      | sync-upstream               | enter   |
      | push-new-branches           | enter   |
      | push-hook                   | enter   |
      | ship-delete-tracking-branch | enter   |
      | sync-before-ship            | enter   |
      | save config file            | enter   |
    Then it runs the commands
      | COMMAND                                        |
      | git config --global alias.append "town append" |
    And global Git setting "alias.append" is now "town append"

  Scenario: don't ask for perennial branches if no branches that could be perennial exist
    Given Git Town is not configured
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS       | DESCRIPTION                                 |
      | welcome                     | enter      |                                             |
      | aliases                     | enter      |                                             |
      | main development branch     | down enter |                                             |
      | perennial branches          |            | no input here since the dialog doesn't show |
      | hosting platform            | enter      |                                             |
      | origin hostname             | enter      |                                             |
      | sync-feature-strategy       | enter      |                                             |
      | sync-perennial-strategy     | enter      |                                             |
      | sync-upstream               | enter      |                                             |
      | push-new-branches           | enter      |                                             |
      | push-hook                   | enter      |                                             |
      | ship-delete-tracking-branch | enter      |                                             |
      | sync-before-ship            | enter      |                                             |
      | save config file            | down enter |                                             |
    Then the main branch is now "main"
    And there are still no perennial branches

  Scenario: remove an existing code hosting override
    Given local Git Town setting "code-hosting-platform" is "github"
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS           | DESCRIPTION                                 |
      | welcome                     | enter          |                                             |
      | aliases                     | enter          |                                             |
      | main development branch     | down enter     |                                             |
      | perennial branches          |                | no input here since the dialog doesn't show |
      | hosting platform            | up up up enter |                                             |
      | origin hostname             | enter          |                                             |
      | sync-feature-strategy       | enter          |                                             |
      | sync-perennial-strategy     | enter          |                                             |
      | sync-upstream               | enter          |                                             |
      | push-new-branches           | enter          |                                             |
      | push-hook                   | enter          |                                             |
      | ship-delete-tracking-branch | enter          |                                             |
      | sync-before-ship            | enter          |                                             |
      | save config file            | down enter     |                                             |
    Then it runs the commands
      | COMMAND                                           |
      | git config --unset git-town.code-hosting-platform |
    And local Git Town setting "code-hosting-platform" is now not set

  Scenario: enter a GitLab token
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main development branch     | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | hosting platform            | up enter          |                                             |
      | gitlab token                | 1 2 3 4 5 6 enter |                                             |
      | origin hostname             | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | sync-before-ship            | enter             |                                             |
      | save config file            | down enter        |                                             |
    Then it runs the commands
      | COMMAND                                          |
      | git config git-town.gitlab-token 123456          |
      | git config git-town.code-hosting-platform gitlab |
    And local Git Town setting "code-hosting-platform" is now "gitlab"
    And local Git Town setting "gitlab-token" is now "123456"

  Scenario: enter a Gitea token
    When I run "git-town config setup" and enter into the dialog:
      | DIALOG                      | KEYS              | DESCRIPTION                                 |
      | welcome                     | enter             |                                             |
      | aliases                     | enter             |                                             |
      | main development branch     | enter             |                                             |
      | perennial branches          |                   | no input here since the dialog doesn't show |
      | hosting platform            | down down enter   |                                             |
      | gitea token                 | 1 2 3 4 5 6 enter |                                             |
      | origin hostname             | enter             |                                             |
      | sync-feature-strategy       | enter             |                                             |
      | sync-perennial-strategy     | enter             |                                             |
      | sync-upstream               | enter             |                                             |
      | push-new-branches           | enter             |                                             |
      | push-hook                   | enter             |                                             |
      | ship-delete-tracking-branch | enter             |                                             |
      | sync-before-ship            | enter             |                                             |
      | save config file            | down enter        |                                             |
    Then it runs the commands
      | COMMAND                                         |
      | git config git-town.gitea-token 123456          |
      | git config git-town.code-hosting-platform gitea |
    And local Git Town setting "code-hosting-platform" is now "gitea"
    And local Git Town setting "gitea-token" is now "123456"
