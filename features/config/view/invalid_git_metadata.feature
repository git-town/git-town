Feature: display invalid Git metadata

  Scenario: invalid sync-feature-strategy
    Given a Git repo with origin
    And local Git setting "git-town.auto-sync" is "zonk"
    And local Git setting "git-town.branch-prefix" is ""
    And local Git setting "git-town.contribution-regex" is "(cont"
    And local Git setting "git-town.detached" is "zonk"
    And local Git setting "git-town.feature-regex" is "(feat"
    And local Git setting "git-town.ignore-uncommitted" is "zonk"
    And local Git setting "git-town.new-branch-type" is "zonk"
    And local Git setting "git-town.observed-regex" is "(obs"
    And local Git setting "git-town.order" is "zonk"
    And local Git setting "git-town.perennial-regex" is "(per"
    And local Git setting "git-town.proposal-breadcrumb" is "zonk"
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
    When I run "git-town config"
    Then Git Town prints:
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
      Ignoring invalid value for "git-town.proposal-breadcrumb": "zonk"
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

      Branches:
        contribution branches: (none)
        contribution regex: (not set)
        feature regex: (not set)
        main branch: main
        observed branches: (none)
        observed regex: (not set)
        parked branches: (none)
        perennial branches: (none)
        perennial regex: (not set)
        prototype branches: (none)
        unknown branch type: feature
        order: asc
        display types: all branch types except "feature" and "main"

      Configuration:
        offline: no
        git user name: user
        git user email: email@example.com

      Create:
        branch prefix: (not set)
        new branch type: (not set)
        share new branches: no
        stash uncommitted changes: yes

      Hosting:
        browser: (not set)
        development remote: origin
        forge type: (not set)
        origin hostname: (not set)
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector: (not set)
        GitHub token: (not set)
        GitLab connector: (not set)
        GitLab token: (not set)

      Propose:
        breadcrumb: none
        breadcrumb direction: down

      Ship:
        delete tracking branch: yes
        ignore uncommitted changes: no
        ship strategy: api

      Sync:
        auto-resolve phantom conflicts: yes
        auto-sync: yes
        run detached: no
        run pre-push hook: yes
        feature sync strategy: merge
        perennial sync strategy: rebase
        prototype sync strategy: merge
        push branches: yes
        sync tags: yes
        sync with upstream: yes
        auto-resolve phantom conflicts: yes
      """
