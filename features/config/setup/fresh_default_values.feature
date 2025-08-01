@messyoutput
Feature: Accepting all default values in a brand-new Git repo leads to a working setup

  Background:
    Given a brand-new Git repo
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                      | KEYS       |
      | welcome                     | enter      |
      | aliases                     | enter      |
      | main branch                 | enter      |
      | perennial branches          |            |
      | perennial regex             | enter      |
      | feature regex               | enter      |
      | contribution regex          | enter      |
      | observed regex              | enter      |
      | new branch type             | enter      |
      | unknown branch type         | enter      |
      | origin hostname             | enter      |
      | forge type                  | enter      |
      | sync feature strategy       | enter      |
      | sync perennial strategy     | enter      |
      | sync prototype strategy     | enter      |
      | sync upstream               | enter      |
      | sync tags                   | enter      |
      | share new branches          | enter      |
      | push hook                   | enter      |
      | ship strategy               | enter      |
      | ship delete tracking branch | enter      |
      | config storage              | down enter |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                         |
      | git config git-town.unknown-branch-type feature |
    And the main branch is still not set
    And there are still no perennial branches
    And local Git setting "git-town.dev-remote" still doesn't exist
    And local Git setting "git-town.new-branch-type" still doesn't exist
    And local Git setting "git-town.main-branch" still doesn't exist
    And local Git setting "git-town.perennial-branches" still doesn't exist
    And local Git setting "git-town.unknown-branch-type" is now "feature"
    And local Git setting "git-town.feature-regex" still doesn't exist
    And local Git setting "git-town.contribution-regex" still doesn't exist
    And local Git setting "git-town.observed-regex" still doesn't exist
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
    And the configuration file is now:
      """
      # More info around this file at https://www.git-town.com/configuration-file

      [branches]
      main = "initial"

      [create]
      new-branch-type = "feature"
      share-new-branches = "no"

      [hosting]
      dev-remote = "origin"

      [ship]
      delete-tracking-branch = true
      strategy = "api"

      [sync]
      feature-strategy = "merge"
      perennial-strategy = "ff-only"
      prototype-strategy = "merge"
      push-hook = true
      tags = true
      upstream = true
      """

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" still doesn't exist
    And global Git setting "alias.diff-parent" still doesn't exist
    And global Git setting "alias.hack" still doesn't exist
    And global Git setting "alias.delete" still doesn't exist
    And global Git setting "alias.prepend" still doesn't exist
    And global Git setting "alias.propose" still doesn't exist
    And global Git setting "alias.rename" still doesn't exist
    And global Git setting "alias.repo" still doesn't exist
    And global Git setting "alias.set-parent" still doesn't exist
    And global Git setting "alias.ship" still doesn't exist
    And global Git setting "alias.sync" still doesn't exist
    And local Git setting "git-town.dev-remote" still doesn't exist
    And local Git setting "git-town.new-branch-type" still doesn't exist
    And local Git setting "git-town.main-branch" still doesn't exist
    And local Git setting "git-town.perennial-branches" still doesn't exist
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.github-token" still doesn't exist
    And local Git setting "git-town.hosting-origin-hostname" still doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.contribution-regex" now doesn't exist
    And local Git setting "git-town.observed-regex" now doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.perennial-regex" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
