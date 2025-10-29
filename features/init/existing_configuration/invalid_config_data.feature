@messyoutput
Feature: Fix invalid configuration data

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And local Git setting "init.defaultbranch" is "main"
    And local Git setting "git-town.feature-regex" is "(feat"
    And local Git setting "git-town.perennial-regex" is "(per"
    And local Git setting "git-town.contribution-regex" is "(cont"
    And local Git setting "git-town.observed-regex" is "(obs"
    And local Git setting "git-town.sync-feature-strategy" is "--help"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                      | KEYS        |
      | welcome                     | enter       |
      | aliases                     | enter       |
      | main branch                 | enter       |
      | origin hostname             | enter       |
      | forge type                  | enter       |
      | enter all                   | down enter  |
      | perennial regex             | p e r enter |
      | feature regex               | enter       |
      | contribution regex          | enter       |
      | observed regex              | enter       |
      | new branch type             | enter       |
      | unknown branch type         | enter       |
      | sync feature strategy       | down enter  |
      | sync perennial strategy     | enter       |
      | sync prototype strategy     | enter       |
      | sync upstream               | enter       |
      | auto-sync                   | enter       |
      | sync-tags                   | enter       |
      | detached                    | enter       |
      | stash                       | enter       |
      | share-new-branches          | enter       |
      | push-branches               | enter       |
      | push-hook                   | enter       |
      | ship-strategy               | enter       |
      | ship-delete-tracking branch | enter       |
      | order                       | enter       |
      | config storage              | enter       |

  @this
  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.auto-sync true                   |
      | git config git-town.detached false                   |
      | git config git-town.new-branch-type feature          |
      | git config git-town.main-branch main                 |
      | git config git-town.perennial-regex per              |
      | git config git-town.unknown-branch-type feature      |
      | git config git-town.order asc                        |
      | git config git-town.push-branches true               |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-strategy api                |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.stash true                       |
      | git config git-town.sync-feature-strategy rebase     |
      | git config git-town.sync-perennial-strategy ff-only  |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-upstream true               |
      | git config git-town.sync-tags true                   |
    And Git Town prints:
      """
      Ignoring invalid value for "git-town.contribution-regex": "(cont"
      Ignoring invalid value for "git-town.feature-regex": "(feat"
      Ignoring invalid value for "git-town.observed-regex": "(obs"
      Ignoring invalid value for "git-town.perennial-regex": "(per"
      Ignoring invalid value for "git-town.sync-feature-strategy": "--help"
      """
    And local Git setting "git-town.sync-feature-strategy" is now "rebase"

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" now doesn't exist
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
    And local Git setting "git-town.sync-feature-strategy" is now "--help"
    And local Git setting "git-town.dev-remote" now doesn't exist
    And local Git setting "git-town.new-branch-type" now doesn't exist
    And local Git setting "git-town.main-branch" now doesn't exist
    And local Git setting "git-town.perennial-branches" now doesn't exist
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.github-token" now doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.contribution-regex" now doesn't exist
    And local Git setting "git-town.observed-regex" now doesn't exist
    And local Git setting "git-town.stash" now doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" now doesn't exist
    And local Git setting "git-town.sync-upstream" now doesn't exist
    And local Git setting "git-town.sync-tags" now doesn't exist
    And local Git setting "git-town.perennial-regex" now doesn't exist
    And local Git setting "git-town.share-new-branches" now doesn't exist
    And local Git setting "git-town.push-hook" now doesn't exist
    And local Git setting "git-town.ship-strategy" now doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" now doesn't exist
