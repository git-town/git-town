@messyoutput
Feature: remove existing configuration in Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE   | LOCATIONS |
      | qa         | (none) | local     |
      | production | (none) | local     |
    And I rename the "origin" remote to "fork"
    And the main branch is "main"
    # keep-sorted start
    And global Git setting "alias.append" is "town append"
    And global Git setting "alias.compress" is "town compress"
    And global Git setting "alias.contribute" is "town contribute"
    And global Git setting "alias.delete" is "town delete"
    And global Git setting "alias.diff-parent" is "town diff-parent"
    And global Git setting "alias.hack" is "town hack"
    And global Git setting "alias.observe" is "town observe"
    And global Git setting "alias.park" is "town park"
    And global Git setting "alias.prepend" is "town prepend"
    And global Git setting "alias.propose" is "town propose"
    And global Git setting "alias.rename" is "town rename"
    And global Git setting "alias.repo" is "town repo"
    And global Git setting "alias.set-parent" is "town set-parent"
    And global Git setting "alias.ship" is "town ship"
    And global Git setting "alias.sync" is "town sync"
    # keep-sorted end
    # keep-sorted start
    And local Git setting "git-town.auto-sync" is "false"
    And local Git setting "git-town.branch-prefix" is "kg-"
    And local Git setting "git-town.contribution-regex" is "other.*"
    And local Git setting "git-town.detached" is "true"
    And local Git setting "git-town.dev-remote" is "fork"
    And local Git setting "git-town.feature-regex" is "user.*"
    And local Git setting "git-town.forge-type" is "github"
    And local Git setting "git-town.hosting-origin-hostname" is "code"
    And local Git setting "git-town.ignore-uncommitted" is "true"
    And local Git setting "git-town.new-branch-type" is "parked"
    And local Git setting "git-town.observed-regex" is "obs.*"
    And local Git setting "git-town.order" is "desc"
    And local Git setting "git-town.perennial-branches" is "qa"
    And local Git setting "git-town.perennial-regex" is "qa.*"
    And local Git setting "git-town.proposals-show-lineage" is "none"
    And local Git setting "git-town.push-branches" is "false"
    And local Git setting "git-town.push-hook" is "false"
    And local Git setting "git-town.share-new-branches" is "push"
    And local Git setting "git-town.ship-delete-tracking-branch" is "false"
    And local Git setting "git-town.ship-strategy" is "squash-merge"
    And local Git setting "git-town.stash" is "false"
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is "ff-only"
    And local Git setting "git-town.sync-prototype-strategy" is "rebase"
    And local Git setting "git-town.sync-tags" is "false"
    And local Git setting "git-town.sync-upstream" is "false"
    And local Git setting "git-town.unknown-branch-type" is "observed"
    # keep-sorted end
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                      | KEYS                                                                        | DESCRIPTION         |
      | welcome                     | enter                                                                       |                     |
      | aliases                     | n enter                                                                     | remove all aliases  |
      | main branch                 | enter                                                                       |                     |
      | perennial branches          | down space enter                                                            |                     |
      | origin hostname             | backspace backspace backspace backspace enter                               | remove the override |
      | forge type                  | up up up up up up enter                                                     | remove the override |
      | enter all                   | down enter                                                                  |                     |
      | perennial regex             | backspace backspace backspace backspace enter                               |                     |
      | feature regex               | backspace backspace backspace backspace backspace backspace enter           |                     |
      | contribution regex          | backspace backspace backspace backspace backspace backspace backspace enter |                     |
      | observed regex              | backspace backspace backspace backspace backspace enter                     |                     |
      | branch prefix               | backspace backspace backspace enter                                         |                     |
      | new branch type             | up enter                                                                    |                     |
      | unknown branch type         | up enter                                                                    |                     |
      | sync feature strategy       | up enter                                                                    |                     |
      | sync perennial strategy     | down enter                                                                  |                     |
      | sync prototype strategy     | up enter                                                                    |                     |
      | sync upstream               | down enter                                                                  |                     |
      | auto sync                   | up enter                                                                    |                     |
      | sync tags                   | down enter                                                                  |                     |
      | detached                    | down enter                                                                  |                     |
      | stash                       | up enter                                                                    |                     |
      | share new branches          | up enter                                                                    | enable              |
      | push branches               | down enter                                                                  | enable              |
      | push hook                   | down enter                                                                  | enable              |
      | ship strategy               | down enter                                                                  |                     |
      | ship delete tracking branch | down enter                                                                  | disable             |
      | ignore-uncommitted          | up enter                                                                    | disable             |
      | order                       | up enter                                                                    |                     |
      | proposals show lineage      | down enter                                                                  |                     |
      | config storage              | enter                                                                       | git metadata        |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config --global --unset alias.append             |
      | git config --global --unset alias.compress           |
      | git config --global --unset alias.contribute         |
      | git config --global --unset alias.diff-parent        |
      | git config --global --unset alias.hack               |
      | git config --global --unset alias.delete             |
      | git config --global --unset alias.observe            |
      | git config --global --unset alias.park               |
      | git config --global --unset alias.prepend            |
      | git config --global --unset alias.propose            |
      | git config --global --unset alias.rename             |
      | git config --global --unset alias.repo               |
      | git config --global --unset alias.set-parent         |
      | git config --global --unset alias.ship               |
      | git config --global --unset alias.sync               |
      | git config git-town.perennial-branches ""            |
      | git config --unset git-town.hosting-origin-hostname  |
      | git config --unset git-town.forge-type               |
      | git config git-town.auto-sync true                   |
      | git config --unset git-town.branch-prefix            |
      | git config --unset git-town.contribution-regex       |
      | git config git-town.detached false                   |
      | git config --unset git-town.feature-regex            |
      | git config git-town.ignore-uncommitted false         |
      | git config git-town.new-branch-type feature          |
      | git config --unset git-town.observed-regex           |
      | git config git-town.order asc                        |
      | git config --unset git-town.perennial-regex          |
      | git config git-town.proposals-show-lineage cli       |
      | git config git-town.push-branches true               |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.ship-strategy api                |
      | git config git-town.stash true                       |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy rebase   |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-tags true                   |
      | git config git-town.sync-upstream true               |
      | git config git-town.unknown-branch-type feature      |
    # keep-sorted start
    And global Git setting "alias.append" now doesn't exist
    And global Git setting "alias.delete" now doesn't exist
    And global Git setting "alias.diff-parent" now doesn't exist
    And global Git setting "alias.hack" now doesn't exist
    And global Git setting "alias.prepend" now doesn't exist
    And global Git setting "alias.propose" now doesn't exist
    And global Git setting "alias.rename" now doesn't exist
    And global Git setting "alias.repo" now doesn't exist
    And global Git setting "alias.set-parent" now doesn't exist
    And global Git setting "alias.ship" now doesn't exist
    And global Git setting "alias.sync" now doesn't exist
    # keep-sorted end
    # keep-sorted start
    And local Git setting "git-town.auto-sync" is now "true"
    And local Git setting "git-town.dev-remote" is now "fork"
    And local Git setting "git-town.ignore-uncommitted" is now "false"
    And local Git setting "git-town.new-branch-type" is now "feature"
    And local Git setting "git-town.order" is now "asc"
    And local Git setting "git-town.proposals-show-lineage" is now "cli"
    And local Git setting "git-town.push-branches" is now "true"
    And local Git setting "git-town.push-hook" is now "true"
    And local Git setting "git-town.share-new-branches" is now "no"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "true"
    And local Git setting "git-town.ship-strategy" is now "api"
    And local Git setting "git-town.stash" is now "true"
    And local Git setting "git-town.sync-feature-strategy" is now "merge"
    And local Git setting "git-town.sync-perennial-strategy" is now "rebase"
    And local Git setting "git-town.sync-tags" is now "true"
    And local Git setting "git-town.sync-upstream" is now "true"
    And local Git setting "git-town.unknown-branch-type" is now "feature"
    And local Git setting "git-town.contribution-regex" now doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.observed-regex" now doesn't exist
    And local Git setting "git-town.perennial-regex" now doesn't exist
    # keep-sorted end
    And the main branch is still "main"
    And there are now no perennial branches

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" is now "town append"
    # keep-sorted start
    And global Git setting "alias.delete" is now "town delete"
    And global Git setting "alias.diff-parent" is now "town diff-parent"
    And global Git setting "alias.hack" is now "town hack"
    And global Git setting "alias.prepend" is now "town prepend"
    And global Git setting "alias.propose" is now "town propose"
    And global Git setting "alias.rename" is now "town rename"
    And global Git setting "alias.repo" is now "town repo"
    And global Git setting "alias.set-parent" is now "town set-parent"
    And global Git setting "alias.ship" is now "town ship"
    And global Git setting "alias.sync" is now "town sync"
    # keep-sorted end
    # keep-sorted start
    And local Git setting "git-town.auto-sync" is now "false"
    And local Git setting "git-town.branch-prefix" is now "kg-"
    And local Git setting "git-town.contribution-regex" is now "other.*"
    And local Git setting "git-town.dev-remote" is now "fork"
    And local Git setting "git-town.feature-regex" is now "user.*"
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.hosting-origin-hostname" is now "code"
    And local Git setting "git-town.ignore-uncommitted" is now "true"
    And local Git setting "git-town.new-branch-type" is now "parked"
    And local Git setting "git-town.observed-regex" is now "obs.*"
    And local Git setting "git-town.order" is now "desc"
    And local Git setting "git-town.perennial-regex" is now "qa.*"
    And local Git setting "git-town.proposals-show-lineage" is now "none"
    And local Git setting "git-town.push-branches" is now "false"
    And local Git setting "git-town.push-hook" is now "false"
    And local Git setting "git-town.share-new-branches" is now "push"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "false"
    And local Git setting "git-town.ship-strategy" is now "squash-merge"
    And local Git setting "git-town.stash" is now "false"
    And local Git setting "git-town.sync-feature-strategy" is now "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is now "ff-only"
    And local Git setting "git-town.sync-prototype-strategy" is now "rebase"
    And local Git setting "git-town.sync-tags" is now "false"
    And local Git setting "git-town.sync-upstream" is now "false"
    And local Git setting "git-town.unknown-branch-type" is now "observed"
    # keep-sorted end
    And the main branch is still "main"
    And the perennial branches are now "qa"
