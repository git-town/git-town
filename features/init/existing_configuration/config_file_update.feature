@messyoutput
Feature: don't ask for information already provided by the config file

  Scenario:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [branches]
      main = "main"
      contribution-regex = "contribute-"
      feature-regex = "feat-"
      observed-regex = "observed-"
      perennial-regex = "release-"
      perennials = ["staging"]
      unknown-type = "observed"
      order = "desc"

      [create]
      branch-prefix = "acme-"
      new-branch-type = "feature"
      share-new-branches = "propose"
      stash = true

      [hosting]
      dev-remote = "something"
      origin-hostname = "github.com"
      forge-type = "github"

      [propose]
      lineage = "none"

      [ship]
      delete-tracking-branch = true
      ignore-uncommitted = true
      strategy = "api"

      [sync]
      auto-sync = false
      detached = false
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      push-branches = false
      push-hook = true
      tags = true
      upstream = true

      [sync-strategy]
      feature-branches = "rebase"
      prototype-branches = "merge"
      perennial-branches = "ff-only"
      """
    And Git Town is not configured
    When I run "git-town init" and enter into the dialogs:
      | DIALOG                | KEYS              |
      | welcome               | enter             |
      | aliases               | enter             |
      | perennial branches    | enter             |
      | github connector type | enter             |
      | github token          | g h - t o k enter |
      | token scope           | enter             |
      | enter all             | down enter        |
      | config storage        | down enter        |
    Then Git Town runs the commands
      | COMMAND                                 |
      | git config git-town.github-token gh-tok |
    And the configuration file is now:
      """
      # See https://www.git-town.com/configuration-file for details

      [branches]
      contribution-regex = "contribute-"
      feature-regex = "feat-"
      main = "main"
      observed-regex = "observed-"
      order = "desc"
      perennials = ["staging"]
      perennial-regex = "release-"
      unknown-type = "observed"

      [create]
      branch-prefix = "acme-"
      new-branch-type = "feature"
      share-new-branches = "propose"
      stash = true

      [hosting]
      dev-remote = "something"
      forge-type = "github"
      github-connector = "api"
      origin-hostname = "github.com"

      [propose]
      lineage = "none"

      [ship]
      delete-tracking-branch = true
      ignore-uncommitted = true
      strategy = "api"

      [sync]
      auto-sync = false
      detached = false
      feature-strategy = "merge"
      perennial-strategy = "rebase"
      prototype-strategy = "merge"
      push-branches = false
      push-hook = true
      tags = true
      upstream = true
      """
