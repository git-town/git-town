Feature: reset invalid configuration

  Scenario: sync-feature-strategy is invalid
    Given a Git repo with origin
    And the main branch is "main"
    # keep-sorted start
    And local Git setting "git-town.auto-sync" is "zonk"
    And local Git setting "git-town.contribution-regex" is "(cont"
    And local Git setting "git-town.detached" is "zonk"
    And local Git setting "git-town.feature-regex" is "(feat"
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
    # keep-sorted end
    When I run "git-town config remove"
    Then Git Town runs no commands
    And Git Town is no longer configured
    # keep-sorted start
    And local Git setting "git-town.auto-sync" now doesn't exist
    And local Git setting "git-town.contribution-regex" now doesn't exist
    And local Git setting "git-town.detached" now doesn't exist
    And local Git setting "git-town.feature-regex" now doesn't exist
    And local Git setting "git-town.new-branch-type" now doesn't exist
    And local Git setting "git-town.observed-regex" now doesn't exist
    And local Git setting "git-town.order" now doesn't exist
    And local Git setting "git-town.perennial-regex" now doesn't exist
    And local Git setting "git-town.push-branches" now doesn't exist
    And local Git setting "git-town.push-hook" now doesn't exist
    And local Git setting "git-town.share-new-branches" now doesn't exist
    And local Git setting "git-town.ship-delete-tracking-branch" now doesn't exist
    And local Git setting "git-town.ship-strategy" now doesn't exist
    And local Git setting "git-town.stash" now doesn't exist
    And local Git setting "git-town.sync-feature-strategy" now doesn't exist
    And local Git setting "git-town.sync-perennial-strategy" now doesn't exist
    And local Git setting "git-town.sync-prototype-strategy" now doesn't exist
    And local Git setting "git-town.sync-tags" now doesn't exist
    And local Git setting "git-town.sync-upstream" now doesn't exist
    And local Git setting "git-town.unknown-branch-type" now doesn't exist
    # keep-sorted end
