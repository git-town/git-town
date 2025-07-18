@messyoutput
Feature: ask for information not provided by the config file

  Scenario:
    Given a Git repo with origin
    And Git Town is not configured
    And the committed configuration file:
      """
      [branches]
      main = "main"

      [hosting]
      dev-remote = "something"
      forge-type = "github"
      origin-hostname = "github.com"

      [ship]
      delete-tracking-branch = false

      [sync]
      tags = false
      upstream = false
      """
    When I run "git-town config setup" and enter into the dialogs:
      | DIALOG                  | KEYS        |
      | welcome                 | enter       |
      | aliases                 | enter       |
      | perennial regex         | 1 1 1 enter |
      | feature regex           | 2 2 2 enter |
      | contribution regex      | 3 3 3 enter |
      | observed regex          | 4 4 4 enter |
      | unknown branch type     | enter       |
      | github connector type   | enter       |
      | github token            | 9 9 9 enter |
      | token scope             | enter       |
      | sync feature strategy   | enter       |
      | sync perennial strategy | enter       |
      | sync prototype strategy | enter       |
      | share new branches      | enter       |
      | push hook               | enter       |
      | new branch type         | enter       |
      | ship strategy           | enter       |
      | config storage          | enter       |
    Then Git Town runs the commands
      | COMMAND                                            |
      | git config git-town.github-token 999               |
      | git config git-town.new-branch-type feature        |
      | git config git-town.github-connector api           |
      | git config git-town.perennial-regex 111            |
      | git config git-town.unknown-branch-type feature    |
      | git config git-town.feature-regex 222              |
      | git config git-town.contribution-regex 333         |
      | git config git-town.observed-regex 444             |
      | git config git-town.push-hook true                 |
      | git config git-town.share-new-branches no          |
      | git config git-town.ship-strategy api              |
      | git config git-town.sync-feature-strategy merge    |
      | git config git-town.sync-perennial-strategy rebase |
      | git config git-town.sync-prototype-strategy merge  |
