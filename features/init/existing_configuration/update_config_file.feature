@messyoutput
Feature: update information in the config file

  @debug @this
  Scenario:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [branches]
      main = "main"
      perennials = ["public"]
      order = "desc"
      
      [create]
      branch-prefix = "acme-"
      
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
      | config storage          | down enter            |
    Then the configuration file is now:
      """
      # See https://www.git-town.com/configuration-file for details
      
      [branches]
      main = "main"
      perennials = ["branch-1"]
      perennial-regex = "peren"
      
      [create]
      new-branch-type = "feature"
      share-new-branches = "no"
      stash = true
      
      [propose]
      lineage = "none"
      
      [ship]
      strategy = "api"
      
      [sync]
      feature-strategy = "merge"
      perennial-strategy = "ff-only"
      prototype-strategy = "merge"
      push-branches = true
      push-hook = true
      """
