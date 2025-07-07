@messyoutput
Feature: remove existing configuration in Git metadata

  Background:
    Given a Git repo with origin
    And I rename the "origin" remote to "fork"
    And the branches
      | NAME       | TYPE   | LOCATIONS |
      | qa         | (none) | local     |
      | production | (none) | local     |
    And the main branch is "main"
    And global Git setting "alias.append" is "town append"
    And global Git setting "alias.diff-parent" is "town diff-parent"
    And global Git setting "alias.hack" is "town hack"
    And global Git setting "alias.delete" is "town delete"
    And global Git setting "alias.prepend" is "town prepend"
    And global Git setting "alias.propose" is "town propose"
    And global Git setting "alias.rename" is "town rename"
    And global Git setting "alias.repo" is "town repo"
    And global Git setting "alias.set-parent" is "town set-parent"
    And global Git setting "alias.ship" is "town ship"
    And global Git setting "alias.sync" is "town sync"
    And local Git setting "git-town.forge-type" is "github"
    And local Git setting "git-town.perennial-branches" is "qa"
    And local Git setting "git-town.perennial-regex" is "qa.*"
    And local Git setting "git-town.feature-regex" is "user.*"
    And local Git setting "git-town.unknown-branch-type" is "observed"
    And local Git setting "git-town.dev-remote" is "fork"
    And local Git setting "git-town.push-hook" is "false"
    And local Git setting "git-town.hosting-origin-hostname" is "code"
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And local Git setting "git-town.sync-prototype-strategy" is "rebase"
    And local Git setting "git-town.sync-upstream" is "false"
    And local Git setting "git-town.sync-tags" is "false"
    And local Git setting "git-town.share-new-branches" is "push"
    And local Git setting "git-town.push-hook" is "false"
    And local Git setting "git-town.new-branch-type" is "parked"
    And local Git setting "git-town.ship-strategy" is "squash-merge"
    And local Git setting "git-town.ship-delete-tracking-branch" is "false"
    When I run "git-town config setup" and enter into the dialogs:
      | DESCRIPTION                             | KEYS                                                              |
      | welcome                                 | enter                                                             |
      | add all aliases                         | n enter                                                           |
      | keep the already configured main branch | enter                                                             |
      | remove the perennial branches           | down space enter                                                  |
      | remove the perennial regex              | backspace backspace backspace backspace enter                     |
      | feature regex                           | backspace backspace backspace backspace backspace backspace enter |
      | unknown branch type                     | up enter                                                          |
      | dev-remote                              | enter                                                             |
      | remove origin hostname                  | backspace backspace backspace backspace enter                     |
      | remove forge type override              | up up up up up enter                                              |
      | sync-feature-strategy                   | up enter                                                          |
      | sync-perennial-strategy                 | down enter                                                        |
      | sync-prototype-strategy                 | up enter                                                          |
      | sync-upstream                           | down enter                                                        |
      | sync-tags                               | down enter                                                        |
      | enable share-new-branches               | up enter                                                          |
      | enable the push hook                    | down enter                                                        |
      | new-branch-type                         | down enter                                                        |
      | change ship-strategy                    | down enter                                                        |
      | disable ship-delete-tracking-branch     | down enter                                                        |
      | save config to Git metadata             | down enter                                                        |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                             |
      | git config --global --unset alias.append            |
      | git config --global --unset alias.diff-parent       |
      | git config --global --unset alias.hack              |
      | git config --global --unset alias.delete            |
      | git config --global --unset alias.prepend           |
      | git config --global --unset alias.propose           |
      | git config --global --unset alias.rename            |
      | git config --global --unset alias.repo              |
      | git config --global --unset alias.set-parent        |
      | git config --global --unset alias.ship              |
      | git config --global --unset alias.sync              |
      | git config --unset git-town.forge-type              |
      | git config --unset git-town.hosting-origin-hostname |
    And global Git setting "alias.append" now doesn't exist
    And global Git setting "alias.diff-parent" now doesn't exist
    And global Git setting "alias.hack" now doesn't exist
    And global Git setting "alias.delete" now doesn't exist
    And global Git setting "alias.prepend" now doesn't exist
    And global Git setting "alias.propose" now doesn't exist
    And global Git setting "alias.rename" now doesn't exist
    And global Git setting "alias.repo" now doesn't exist
    And global Git setting "alias.set-parent" now doesn't exist
    And global Git setting "alias.ship" now doesn't exist
    And global Git setting "alias.sync" now doesn't exist
    And the main branch is still "main"
    And the perennial branches are now "production"
    And local Git setting "git-town.dev-remote" is now "fork"
    And local Git setting "git-town.new-branch-type" is now "prototype"
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.sync-feature-strategy" is now "compress"
    And local Git setting "git-town.sync-perennial-strategy" is now "ff-only"
    And local Git setting "git-town.sync-upstream" is now "false"
    And local Git setting "git-town.sync-tags" is now "false"
    And local Git setting "git-town.perennial-regex" now doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.unknown-branch-type" is now "parked"
    And local Git setting "git-town.share-new-branches" is now "no"
    And local Git setting "git-town.push-hook" is now "false"
    And local Git setting "git-town.ship-strategy" is now "api"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "true"

  Scenario: undo
    When I run "git-town undo"
    Then the main branch is still "main"
    And the perennial branches are now "qa"
    And global Git setting "alias.append" is now "town append"
    And global Git setting "alias.diff-parent" is now "town diff-parent"
    And global Git setting "alias.hack" is now "town hack"
    And global Git setting "alias.delete" is now "town delete"
    And global Git setting "alias.prepend" is now "town prepend"
    And global Git setting "alias.propose" is now "town propose"
    And global Git setting "alias.rename" is now "town rename"
    And global Git setting "alias.repo" is now "town repo"
    And global Git setting "alias.set-parent" is now "town set-parent"
    And global Git setting "alias.ship" is now "town ship"
    And global Git setting "alias.sync" is now "town sync"
    And local Git setting "git-town.dev-remote" is now "fork"
    And local Git setting "git-town.new-branch-type" is now "parked"
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.perennial-regex" is now "qa.*"
    And local Git setting "git-town.feature-regex" is now "user.*"
    And local Git setting "git-town.unknown-branch-type" is now "observed"
    And local Git setting "git-town.share-new-branches" is now "push"
    And local Git setting "git-town.push-hook" is now "true"
    And local Git setting "git-town.hosting-origin-hostname" is now "code"
    And local Git setting "git-town.sync-feature-strategy" is now "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is now "rebase"
    And local Git setting "git-town.sync-upstream" is now "true"
    And local Git setting "git-town.sync-tags" is now "true"
    And local Git setting "git-town.share-new-branches" is now "push"
    And local Git setting "git-town.push-hook" is now "true"
    And local Git setting "git-town.ship-strategy" is now "squash-merge"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "false"
