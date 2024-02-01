Feature: migrate existing configuration in Git metadata to a config file

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
      | save config to config file                | enter                  |

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
    And the main branch is now not set
    And there are now no perennial branches
    And local Git Town setting "code-hosting-platform" no longer exists
    And local Git Town setting "github-token" is now "123456"
    And local Git Town setting "code-hosting-origin-hostname" no longer exists
    And local Git Town setting "sync-feature-strategy" no longer exists
    And local Git Town setting "sync-perennial-strategy" no longer exists
    And local Git Town setting "sync-upstream" no longer exists
    And local Git Town setting "push-new-branches" no longer exists
    And local Git Town setting "push-hook" no longer exists
    And local Git Town setting "ship-delete-tracking-branch" no longer exists
    And local Git Town setting "sync-before-ship" no longer exists
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
      push-new-branches = true

      # Should "git ship" delete the tracking branch?
      # You want to disable this if your code hosting system
      # (GitHub, GitLab, etc) deletes head branches when
      # merging pull requests through its UI.
      ship-delete-tracking-branch = false

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
      sync-before-ship = true

      # Should "git sync" also fetch updates from the upstream remote?
      #
      # If an "upstream" remote exists, and this setting is enabled,
      # "git sync" will also update the local main branch
      # with commits from the main branch at the upstream remote.
      #
      # This is useful if the repository you work on is a fork,
      # and you want to keep it in sync with the repo it was forked from.
      sync-upstream = false

      [branches]

        # The main branch is the branch from which you cut new feature branches,
        # and into which you ship feature branches when they are done.
        # This branch is often called "main", "master", or "development".
        main = "main"

        # Perennial branches are long-lived branches.
        # They are never shipped and have no ancestors.
        # Typically, perennial branches have names like
        # "development", "staging", "qa", "production", etc.
        perennials = ["production"]

      [sync-strategy]

        # How should Git Town synchronize feature branches?
        # Feature branches are short-lived branches cut from
        # the main branch and shipped back into the main branch.
        # Typically you develop features and bug fixes on them,
        # hence their name.
        feature-branches = "rebase"

        # How should Git Town synchronize perennial branches?
        # Perennial branches have no parent branch.
        # The only updates they receive are additional commits
        # made to their tracking branch somewhere else.
        perennial-branches = "merge"
      """
