@messyoutput
Feature: don't change existing extended information when changing basic information

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | LOCATIONS     |
      | qa         | perennial | local, origin |
      | production | (none)    | local, origin |
    And the main branch is "main"
    And local Git setting "git-town.new-branch-type" is "parked"
    And local Git setting "git-town.share-new-branches" is "no"
    And local Git setting "git-town.order" is "desc"
    And local Git setting "git-town.push-branches" is "false"
    And local Git setting "git-town.push-hook" is "false"
    And local Git setting "git-town.auto-sync" is "false"
    And local Git setting "git-town.sync-tags" is "false"
    And local Git setting "git-town.detached" is "false"
    And local Git setting "git-town.sync-feature-strategy" is "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is "rebase"
    And local Git setting "git-town.sync-prototype-strategy" is "rebase"
    And local Git setting "git-town.sync-upstream" is "false"
    And local Git setting "git-town.ship-delete-tracking-branch" is "false"
    And local Git setting "git-town.perennial-regex" is "per"
    And local Git setting "git-town.feature-regex" is "feat"
    And local Git setting "git-town.contribution-regex" is "cont"
    And local Git setting "git-town.observed-regex" is "obs"
    And local Git setting "git-town.unknown-branch-type" is "observed"
    And local Git setting "git-town.share-new-branches" is "push"
    And local Git setting "git-town.push-branches" is "true"
    And local Git setting "git-town.push-hook" is "true"
    And local Git setting "git-town.ship-strategy" is "fast-forward"
    And local Git setting "git-town.ship-delete-tracking-branch" is "true"
    And local Git setting "git-town.proposals-show-lineage" is "cli"
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                | KEYS                   |
      | welcome               | enter                  |
      | aliases               | enter                  |
      | main branch           | enter                  |
      | perennial branches    | space down space enter |
      | origin hostname       | c o d e enter          |
      | forge type            | up up enter            |
      | github connector type | enter                  |
      | github token          | g h - t o k enter      |
      | token scope           | enter                  |
      | enter all             | enter                  |
      | config storage        | enter                  |

  Scenario: result
    Then Git Town runs the commands
      | COMMAND                                                |
      | git config git-town.github-token gh-tok                |
      | git config git-town.perennial-branches "production qa" |
      | git config git-town.hosting-origin-hostname code       |
      | git config git-town.forge-type github                  |
    And local Git setting "git-town.auto-sync" is still "false"
    And local Git setting "git-town.detached" is still "false"
    And local Git setting "git-town.perennial-branches" is now "production qa"
    And local Git setting "git-town.new-branch-type" is still "parked"
    And local Git setting "git-town.forge-type" is now "github"
    And local Git setting "git-town.github-token" is now "gh-tok"
    And local Git setting "git-town.hosting-origin-hostname" is now "code"
    And local Git setting "git-town.sync-feature-strategy" is still "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is still "rebase"
    And local Git setting "git-town.sync-prototype-strategy" is still "rebase"
    And local Git setting "git-town.sync-upstream" is still "false"
    And local Git setting "git-town.sync-tags" is still "false"
    And local Git setting "git-town.perennial-regex" is still "per"
    And local Git setting "git-town.feature-regex" is still "feat"
    And local Git setting "git-town.contribution-regex" is still "cont"
    And local Git setting "git-town.observed-regex" is still "obs"
    And local Git setting "git-town.unknown-branch-type" is still "observed"
    And local Git setting "git-town.share-new-branches" is still "push"
    And local Git setting "git-town.push-branches" is still "true"
    And local Git setting "git-town.push-hook" is still "true"
    And local Git setting "git-town.ship-strategy" is still "fast-forward"
    And local Git setting "git-town.ship-delete-tracking-branch" is still "true"
    And local Git setting "git-town.proposals-show-lineage" is still "cli"
    And local Git setting "git-town.stash" still doesn't exist
    And local Git setting "git-town.dev-remote" still doesn't exist
    And the main branch is now "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And global Git setting "alias.append" still doesn't exist
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
    And the main branch is now "main"
    And the perennial branches are now "qa"
    And local Git setting "git-town.hosting-origin-hostname" now doesn't exist
    And local Git setting "git-town.forge-type" now doesn't exist
    And local Git setting "git-town.auto-sync" is still "false"
    And local Git setting "git-town.detached" is still "false"
    And local Git setting "git-town.new-branch-type" is still "parked"
    And local Git setting "git-town.share-new-branches" is still "push"
    And local Git setting "git-town.push-branches" is still "true"
    And local Git setting "git-town.push-hook" is still "true"
    And local Git setting "git-town.ship-delete-tracking-branch" is still "true"
    And local Git setting "git-town.sync-feature-strategy" is still "rebase"
    And local Git setting "git-town.sync-perennial-strategy" is still "rebase"
    And local Git setting "git-town.sync-prototype-strategy" is still "rebase"
    And local Git setting "git-town.sync-upstream" is still "false"
    And local Git setting "git-town.perennial-regex" is still "per"
    And local Git setting "git-town.feature-regex" is still "feat"
    And local Git setting "git-town.contribution-regex" is still "cont"
    And local Git setting "git-town.observed-regex" is still "obs"
    And local Git setting "git-town.unknown-branch-type" is still "observed"
    And local Git setting "git-town.ship-strategy" is still "fast-forward"
    And local Git setting "git-town.proposals-show-lineage" is still "cli"
    And local Git setting "git-town.github-token" now doesn't exist
    And local Git setting "git-town.stash" still doesn't exist
