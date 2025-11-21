@messyoutput
Feature: ask for information not provided by the config file

  Scenario:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [branches]
      main = "main"
      perennials = ["public"]
      order = "desc"
      
      [hosting]
      dev-remote = "something"
      forge-type = "github"
      origin-hostname = "github.com"
      
      [ship]
      delete-tracking-branch = false
      
      [sync]
      auto-sync = false
      tags = false
      upstream = false
      """
    And the branches
      | NAME     | TYPE   | LOCATIONS     |
      | branch-1 | (none) | local, origin |
    And Git Town is not configured
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                  | KEYS                  |
      | welcome                 | enter                 |
      | aliases                 | enter                 |
      | perennial branches      | space enter           |
      | github connector type   | enter                 |
      | github token            | g h - t o k e n enter |
      | token scope             | enter                 |
      | enter all               | down enter            |
      | perennial regex         | p e r e n enter       |
      | feature regex           | f e a t enter         |
      | contribution regex      | c o n t enter         |
      | observed regex          | o b s enter           |
      | branch prefix           | a c m e - enter       |
      | new branch type         | enter                 |
      | unknown branch type     | enter                 |
      | sync feature strategy   | enter                 |
      | sync perennial strategy | enter                 |
      | sync prototype strategy | enter                 |
      | detached                | enter                 |
      | stash                   | enter                 |
      | share new branches      | enter                 |
      | push branches           | enter                 |
      | push hook               | enter                 |
      | ship strategy           | enter                 |
      | proposals show lineage  | enter                 |
      | config storage          | enter                 |
    Then Git Town runs the commands
      | COMMAND                                             |
      | git config git-town.github-token gh-token           |
      | git config git-town.perennial-branches branch-1     |
      | git config git-town.github-connector api            |
      | git config git-town.branch-prefix acme-             |
      | git config git-town.detached false                  |
      | git config git-town.new-branch-type feature         |
      | git config git-town.perennial-regex peren           |
      | git config git-town.unknown-branch-type feature     |
      | git config git-town.feature-regex feat              |
      | git config git-town.contribution-regex cont         |
      | git config git-town.observed-regex obs              |
      | git config git-town.proposals-show-lineage none     |
      | git config git-town.push-branches true              |
      | git config git-town.push-hook true                  |
      | git config git-town.share-new-branches no           |
      | git config git-town.ship-strategy api               |
      | git config git-town.stash true                      |
      | git config git-town.sync-feature-strategy merge     |
      | git config git-town.sync-perennial-strategy ff-only |
      | git config git-town.sync-prototype-strategy merge   |
