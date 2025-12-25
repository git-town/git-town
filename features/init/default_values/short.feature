@messyoutput
Feature: Accepting all default values leads to a working setup

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And local Git setting "init.defaultbranch" is "main"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG          | KEYS       |
      | welcome         | enter      |
      | aliases         | enter      |
      | main branch     | enter      |
      | origin hostname | enter      |
      | forge type      | enter      |
      | enter all       | enter      |
      | config storage  | down enter |

  Scenario: result
    Then Git Town runs no commands
    # keep-sorted start
    And local Git setting "git-town.contribution-regex" still doesn't exist
    And local Git setting "git-town.dev-remote" still doesn't exist
    And local Git setting "git-town.feature-regex" still doesn't exist
    And local Git setting "git-town.forge-type" still doesn't exist
    And local Git setting "git-town.main-branch" still doesn't exist
    And local Git setting "git-town.new-branch-type" still doesn't exist
    And local Git setting "git-town.observed-regex" still doesn't exist
    And local Git setting "git-town.perennial-branches" still doesn't exist
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.stash" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
    And local Git setting "git-town.unknown-branch-type" still doesn't exist
    # keep-sorted end
    And the configuration file is now:
      """
      # See https://www.git-town.com/configuration-file for details

      [branches]
      main = "main"
      """
    And the main branch is still not set
    And there are still no perennial branches

  Scenario: undo
    When I run "git-town undo"
    Then global Git setting "alias.append" still doesn't exist
    # keep-sorted start
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
    # keep-sorted end
    # keep-sorted start
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
    And local Git setting "git-town.push-hook" still doesn't exist
    And local Git setting "git-town.share-new-branches" still doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" still doesn't exist
    And local Git setting "git-town.ship-strategy" still doesn't exist
    And local Git setting "git-town.stash" still doesn't exist
    And local Git setting "git-town.sync-feature-strategy" still doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" still doesn't exist
    And local Git setting "git-town.sync-tags" still doesn't exist
    And local Git setting "git-town.sync-upstream" still doesn't exist
  # keep-sorted end
