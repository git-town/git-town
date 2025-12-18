@messyoutput
Feature: change existing information in Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | LOCATIONS     |
      | qa         | perennial | local, origin |
      | production | (none)    | local, origin |
    And the main branch is "main"
    # keep-sorted start
    And local Git setting "git-town.ignore-uncommitted" is "true"
    And local Git setting "git-town.new-branch-type" is "parked"
    And local Git setting "git-town.share-new-branches" is "no"
    And local Git setting "git-town.order" is "desc"
    And local Git setting "git-town.push-branches" is "false"
    And local Git setting "git-town.push-hook" is "false"
    And local Git setting "git-town.auto-sync" is "false"
    And local Git setting "git-town.sync-tags" is "false"
    And local Git setting "git-town.detached" is "false"
    And local Git setting "git-town.ship-delete-tracking-branch" is "false"
    # keep-sorted end
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                      | KEYS                   |
      | welcome                     | enter                  |
      | aliases                     | a enter                |
      | main branch                 | enter                  |
      | perennial branches          | space down space enter |
      | origin hostname             | c o d e enter          |
      | forge type                  | up up enter            |
      | github connector type       | enter                  |
      | github token                | g h - t o k enter      |
      | token scope                 | enter                  |
      | enter all                   | down enter             |
      | perennial regex             | p e r enter            |
      | feature regex               | f e a t enter          |
      | contribution regex          | c o n t enter          |
      | observed regex              | o b s enter            |
      | branch prefix               | k g - enter            |
      | new branch type             | down enter             |
      | unknown branch type         | down enter             |
      | sync feature strategy       | down enter             |
      | sync perennial strategy     | down enter             |
      | sync prototype strategy     | down enter             |
      | sync upstream               | down enter             |
      | auto sync                   | down enter             |
      | sync tags                   | down enter             |
      | detached                    | up enter               |
      | stash                       | down enter             |
      | share-new-branches          | down enter             |
      | push branches               | down enter             |
      | push hook                   | down enter             |
      | ship strategy               | down down enter        |
      | ship delete tracking branch | down enter             |
      | ignore-uncommitted          | up enter               |
      | order                       | up enter               |
      | proposals show lineage      | down enter             |
      | config storage              | enter                  |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                                  |
      | git config --global alias.append "town append"           |
      | git config --global alias.compress "town compress"       |
      | git config --global alias.contribute "town contribute"   |
      | git config --global alias.diff-parent "town diff-parent" |
      | git config --global alias.down "town down"               |
      | git config --global alias.hack "town hack"               |
      | git config --global alias.delete "town delete"           |
      | git config --global alias.observe "town observe"         |
      | git config --global alias.park "town park"               |
      | git config --global alias.prepend "town prepend"         |
      | git config --global alias.propose "town propose"         |
      | git config --global alias.rename "town rename"           |
      | git config --global alias.repo "town repo"               |
      | git config --global alias.set-parent "town set-parent"   |
      | git config --global alias.ship "town ship"               |
      | git config --global alias.sync "town sync"               |
      | git config --global alias.up "town up"                   |
      | git config git-town.github-token gh-tok                  |
      | git config git-town.perennial-branches "production qa"   |
      | git config git-town.hosting-origin-hostname code         |
      | git config git-town.forge-type github                    |
      | git config git-town.github-connector api                 |
      | git config git-town.auto-sync true                       |
      | git config git-town.branch-prefix kg-                    |
      | git config git-town.detached true                        |
      | git config git-town.new-branch-type prototype            |
      | git config git-town.perennial-regex per                  |
      | git config git-town.unknown-branch-type observed         |
      | git config git-town.feature-regex feat                   |
      | git config git-town.contribution-regex cont              |
      | git config git-town.observed-regex obs                   |
      | git config git-town.order asc                            |
      | git config git-town.proposals-show-lineage cli           |
      | git config git-town.push-branches true                   |
      | git config git-town.push-hook true                       |
      | git config git-town.share-new-branches push              |
      | git config git-town.ship-strategy fast-forward           |
      | git config git-town.ship-delete-tracking-branch true     |
      | git config git-town.ignore-uncommitted false             |
      | git config git-town.stash false                          |
      | git config git-town.sync-feature-strategy rebase         |
      | git config git-town.sync-perennial-strategy rebase       |
      | git config git-town.sync-prototype-strategy rebase       |
      | git config git-town.sync-upstream false                  |
      | git config git-town.sync-tags true                       |
    # keep-sorted start
    And global Git setting "alias.append" is now "town append"
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
    And local Git setting "git-town.contribution-regex" is now "cont"
    And local Git setting "git-town.feature-regex" is now "feat"
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.github-token" is now "gh-tok"
    And local Git setting "git-town.hosting-origin-hostname" is now "code"
    And local Git setting "git-town.ignore-uncommitted" is now "false"
    And local Git setting "git-town.new-branch-type" is now "prototype"
    And local Git setting "git-town.observed-regex" is now "obs"
    And local Git setting "git-town.perennial-branches" is now "production qa"
    And local Git setting "git-town.perennial-regex" is now "per"
    And local Git setting "git-town.proposals-show-lineage" is now "cli"
    And local Git setting "git-town.push-branches" is now "true"
    And local Git setting "git-town.push-hook" is now "true"
    And local Git setting "git-town.share-new-branches" is now "push"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "true"
    And local Git setting "git-town.ship-strategy" is now "fast-forward"
    And local Git setting "git-town.stash" is now "false"
    And local Git setting "git-town.sync-feature-strategy" is now "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is now "rebase"
    And local Git setting "git-town.sync-prototype-strategy" is now "rebase"
    And local Git setting "git-town.sync-tags" is now "true"
    And local Git setting "git-town.sync-upstream" is now "false"
    And local Git setting "git-town.unknown-branch-type" is now "observed"
    And local Git setting "git-town.dev-remote" still doesn't exist
    # keep-sorted end
    And the main branch is now "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
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
    And local Git setting "git-town.auto-sync" is now "false"
    And local Git setting "git-town.ignore-uncommitted" is now "true"
    And local Git setting "git-town.new-branch-type" is now "parked"
    And local Git setting "git-town.push-branches" is now "false"
    And local Git setting "git-town.push-hook" is now "false"
    And local Git setting "git-town.share-new-branches" is now "no"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "false"
    And local Git setting "git-town.contribution-regex" now doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.observed-regex" now doesn't exist
    And local Git setting "git-town.perennial-regex" now doesn't exist
    And local Git setting "git-town.proposals-show-lineage" now doesn't exist
    And local Git setting "git-town.ship-strategy" now doesn't exist
    And local Git setting "git-town.stash" now doesn't exist
    And local Git setting "git-town.sync-feature-strategy" now doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" now doesn't exist
    And local Git setting "git-town.sync-prototype-strategy" now doesn't exist
    And local Git setting "git-town.sync-upstream" now doesn't exist
    And local Git setting "git-town.unknown-branch-type" now doesn't exist
    # keep-sorted end
    And the main branch is now "main"
    And the perennial branches are now "qa"
