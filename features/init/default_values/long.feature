@messyoutput
Feature: Accepting all default values leads to a working setup

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE   | LOCATIONS     |
      | dev        | (none) | local, origin |
      | production | (none) | local, origin |
    And Git Town is not configured
    And local Git setting "init.defaultbranch" is "main"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                            | KEYS       |
      | welcome                           | enter      |
      | aliases                           | enter      |
      | main branch                       | enter      |
      | perennial branches                | enter      |
      | origin hostname                   | enter      |
      | forge type                        | enter      |
      | enter all                         | down enter |
      | perennial regex                   | enter      |
      | feature regex                     | enter      |
      | contribution regex                | enter      |
      | observed regex                    | enter      |
      | branch prefix                     | enter      |
      | new branch type                   | enter      |
      | unknown branch type               | enter      |
      | sync feature strategy             | enter      |
      | sync perennial strategy           | enter      |
      | sync prototype strategy           | enter      |
      | sync upstream                     | enter      |
      | auto-sync                         | enter      |
      | sync-tags                         | enter      |
      | detached                          | enter      |
      | stash                             | enter      |
      | share-new-branches                | enter      |
      | push-branches                     | enter      |
      | push-hook                         | enter      |
      | ship-strategy                     | enter      |
      | ship-delete-tracking branch       | enter      |
      | ignore-uncommitted                | enter      |
      | order                             | enter      |
      | proposal breadcrumb               | enter      |
      | proposals breadcrumb single stack | enter      |
      | config storage                    | enter      |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                              |
      | git config git-town.main-branch main                 |
      | git config git-town.auto-sync true                   |
      | git config git-town.detached false                   |
      | git config git-town.ignore-uncommitted true          |
      | git config git-town.new-branch-type feature          |
      | git config git-town.order asc                        |
      | git config git-town.proposal-breadcrumb none         |
      | git config git-town.proposal-breadcrumb-single true  |
      | git config git-town.push-branches true               |
      | git config git-town.push-hook true                   |
      | git config git-town.share-new-branches no            |
      | git config git-town.ship-delete-tracking-branch true |
      | git config git-town.ship-strategy api                |
      | git config git-town.stash true                       |
      | git config git-town.sync-feature-strategy merge      |
      | git config git-town.sync-perennial-strategy ff-only  |
      | git config git-town.sync-prototype-strategy merge    |
      | git config git-town.sync-tags true                   |
      | git config git-town.sync-upstream true               |
      | git config git-town.unknown-branch-type feature      |
    And there are still no perennial branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And global Git setting "alias.append" still doesn't exist
    And global Git setting "alias.delete" still doesn't exist
    And global Git setting "alias.diff-parent" still doesn't exist
    And global Git setting "alias.hack" still doesn't exist
    And global Git setting "alias.prepend" still doesn't exist
    And global Git setting "alias.propose" still doesn't exist
    And global Git setting "alias.rename" still doesn't exist
    And global Git setting "alias.repo" still doesn't exist
    And global Git setting "alias.set-parent" still doesn't exist
    And global Git setting "alias.ship" still doesn't exist
    And global Git setting "alias.sync" still doesn't exist
    And local Git setting "git-town.branch-prefix" still doesn't exist
    And local Git setting "git-town.contribution-regex" now doesn't exist
    And local Git setting "git-town.dev-remote" still doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" still doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" still doesn't exist
    And local Git setting "git-town.main-branch" still doesn't exist
    And local Git setting "git-town.new-branch-type" still doesn't exist
    And local Git setting "git-town.observed-regex" now doesn't exist
    And local Git setting "git-town.perennial-branches" still doesn't exist
    And local Git setting "git-town.perennial-regex" still doesn't exist
    And local Git setting "git-town.proposal-breadcrumb" still doesn't exist
    And local Git setting "git-town.proposal-breadcrumb-single" still doesn't exist
    And local Git setting "git-town.proposals-show-lineage" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.stash" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
