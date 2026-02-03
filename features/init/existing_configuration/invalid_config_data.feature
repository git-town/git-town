@messyoutput
Feature: Fix invalid configuration data

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And local Git setting "git-town.auto-sync" is "zonk"
    And local Git setting "git-town.branch-prefix" is "xx"
    And local Git setting "git-town.contribution-regex" is "(cont"
    And local Git setting "git-town.detached" is "zonk"
    And local Git setting "git-town.feature-regex" is "(feat"
    And local Git setting "git-town.ignore-uncommitted" is "zonk"
    And local Git setting "git-town.new-branch-type" is "zonk"
    And local Git setting "git-town.observed-regex" is "(obs"
    And local Git setting "git-town.order" is "zonk"
    And local Git setting "git-town.perennial-regex" is "(per"
    And local Git setting "git-town.push-branches" is "zonk"
    And local Git setting "git-town.push-hook" is "zonk"
    And local Git setting "git-town.share-new-branches" is "zonk"
    And local Git setting "git-town.ship-delete-tracking-branch" is "zonk"
    And local Git setting "git-town.ship-strategy" is "zonk"
    And local Git setting "git-town.stash" is "zonk"
    And local Git setting "git-town.sync-feature-strategy" is "--help"
    And local Git setting "git-town.sync-perennial-strategy" is "zonk"
    And local Git setting "git-town.sync-prototype-strategy" is "zonk"
    And local Git setting "git-town.sync-tags" is "zonk"
    And local Git setting "git-town.sync-upstream" is "zonk"
    And local Git setting "git-town.unknown-branch-type" is "zonk"
    And local Git setting "init.defaultbranch" is "main"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                      | KEYS                          |
      | welcome                     | enter                         |
      | aliases                     | enter                         |
      | main branch                 | enter                         |
      | origin hostname             | enter                         |
      | forge type                  | enter                         |
      | enter all                   | down enter                    |
      | perennial regex             | p e r enter                   |
      | feature regex               | f e a t enter                 |
      | contribution regex          | c o n t enter                 |
      | observed regex              | o b s enter                   |
      | branch prefix               | backspace backspace a b enter |
      | new branch type             | down enter                    |
      | unknown branch type         | down enter                    |
      | sync feature strategy       | down enter                    |
      | sync perennial strategy     | down enter                    |
      | sync prototype strategy     | down enter                    |
      | sync upstream               | down enter                    |
      | auto-sync                   | down enter                    |
      | sync-tags                   | down enter                    |
      | detached                    | down enter                    |
      | stash                       | down enter                    |
      | share-new-branches          | down enter                    |
      | push-branches               | down enter                    |
      | push-hook                   | down enter                    |
      | ship-strategy               | down enter                    |
      | ship-delete-tracking branch | down enter                    |
      | ignore-uncommitted          | down enter                    |
      | order                       | down enter                    |
      | proposal breadcrumb         | enter                         |
      | config storage              | enter                         |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                               |
      | git config git-town.main-branch main                  |
      | git config git-town.auto-sync false                   |
      | git config git-town.branch-prefix ab                  |
      | git config git-town.contribution-regex cont           |
      | git config git-town.detached true                     |
      | git config git-town.feature-regex feat                |
      | git config git-town.ignore-uncommitted false          |
      | git config git-town.new-branch-type parked            |
      | git config git-town.observed-regex obs                |
      | git config git-town.order desc                        |
      | git config git-town.perennial-regex per               |
      | git config git-town.proposal-breadcrumb none          |
      | git config git-town.push-branches false               |
      | git config git-town.push-hook false                   |
      | git config git-town.share-new-branches push           |
      | git config git-town.ship-delete-tracking-branch false |
      | git config git-town.ship-strategy always-merge        |
      | git config git-town.stash false                       |
      | git config git-town.sync-feature-strategy rebase      |
      | git config git-town.sync-perennial-strategy rebase    |
      | git config git-town.sync-prototype-strategy rebase    |
      | git config git-town.sync-tags false                   |
      | git config git-town.sync-upstream false               |
      | git config git-town.unknown-branch-type observed      |
    And Git Town prints:
      """
      Ignoring invalid value for "git-town.auto-sync": "zonk"
      Ignoring invalid value for "git-town.contribution-regex": "(cont"
      Ignoring invalid value for "git-town.detached": "zonk"
      Ignoring invalid value for "git-town.feature-regex": "(feat"
      Ignoring invalid value for "git-town.ignore-uncommitted": "zonk"
      Ignoring invalid value for "git-town.new-branch-type": "zonk"
      Ignoring invalid value for "git-town.observed-regex": "(obs"
      Ignoring invalid value for "git-town.order": "zonk"
      Ignoring invalid value for "git-town.perennial-regex": "(per"
      Ignoring invalid value for "git-town.push-branches": "zonk"
      Ignoring invalid value for "git-town.push-hook": "zonk"
      Ignoring invalid value for "git-town.share-new-branches": "zonk"
      Ignoring invalid value for "git-town.ship-delete-tracking-branch": "zonk"
      Ignoring invalid value for "git-town.ship-strategy": "zonk"
      Ignoring invalid value for "git-town.stash": "zonk"
      Ignoring invalid value for "git-town.sync-feature-strategy": "--help"
      Ignoring invalid value for "git-town.sync-perennial-strategy": "zonk"
      Ignoring invalid value for "git-town.sync-prototype-strategy": "zonk"
      Ignoring invalid value for "git-town.sync-tags": "zonk"
      Ignoring invalid value for "git-town.sync-upstream": "zonk"
      Ignoring invalid value for "git-town.unknown-branch-type": "zonk"
      """
    And local Git setting "git-town.auto-sync" is now "false"
    And local Git setting "git-town.contribution-regex" is now "cont"
    And local Git setting "git-town.detached" is now "true"
    And local Git setting "git-town.feature-regex" is now "feat"
    And local Git setting "git-town.new-branch-type" is now "parked"
    And local Git setting "git-town.observed-regex" is now "obs"
    And local Git setting "git-town.order" is now "desc"
    And local Git setting "git-town.perennial-regex" is now "per"
    And local Git setting "git-town.push-branches" is now "false"
    And local Git setting "git-town.push-hook" is now "false"
    And local Git setting "git-town.share-new-branches" is now "push"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "false"
    And local Git setting "git-town.ship-strategy" is now "always-merge"
    And local Git setting "git-town.stash" is now "false"
    And local Git setting "git-town.sync-feature-strategy" is now "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is now "rebase"
    And local Git setting "git-town.sync-prototype-strategy" is now "rebase"
    And local Git setting "git-town.sync-tags" is now "false"
    And local Git setting "git-town.sync-upstream" is now "false"
    And local Git setting "git-town.unknown-branch-type" is now "observed"

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" now doesn't exist
    And global Git setting "alias.compress" now doesn't exist
    And global Git setting "alias.contribute" now doesn't exist
    And global Git setting "alias.delete" now doesn't exist
    And global Git setting "alias.diff-parent" now doesn't exist
    And global Git setting "alias.down" now doesn't exist
    And global Git setting "alias.hack" now doesn't exist
    And global Git setting "alias.observe" now doesn't exist
    And global Git setting "alias.park" now doesn't exist
    And global Git setting "alias.prepend" now doesn't exist
    And global Git setting "alias.propose" now doesn't exist
    And global Git setting "alias.rename" now doesn't exist
    And global Git setting "alias.repo" now doesn't exist
    And global Git setting "alias.set-parent" now doesn't exist
    And global Git setting "alias.ship" now doesn't exist
    And global Git setting "alias.sync" now doesn't exist
    And global Git setting "alias.up" now doesn't exist
    And local Git setting "git-town.contribution-regex" is now "(cont"
    And local Git setting "git-town.dev-remote" now doesn't exist
    And local Git setting "git-town.feature-regex" is now "(feat"
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.main-branch" now doesn't exist
    And local Git setting "git-town.new-branch-type" is now "zonk"
    And local Git setting "git-town.observed-regex" is now "(obs"
    And local Git setting "git-town.perennial-branches" now doesn't exist
    And local Git setting "git-town.perennial-regex" is now "(per"
    And local Git setting "git-town.push-hook" is now "zonk"
    And local Git setting "git-town.share-new-branches" is now "zonk"
    And local Git setting "git-town.ship-delete-tracking-branch" is now "zonk"
    And local Git setting "git-town.ship-strategy" is now "zonk"
    And local Git setting "git-town.stash" is now "zonk"
    And local Git setting "git-town.sync-feature-strategy" is now "--help"
    And local Git setting "git-town.sync-perennial-strategy" is now "zonk"
    And local Git setting "git-town.sync-tags" is now "zonk"
    And local Git setting "git-town.sync-upstream" is now "zonk"
