Feature: display information from config file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS     |
      | contribution-1 | contribution |        | local, origin |
      | contribution-2 | contribution |        | local, origin |
      | observed-1     | observed     |        | local, origin |
      | observed-2     | observed     |        | local, origin |
      | perennial-1    | perennial    |        | local         |
      | perennial-2    | perennial    |        | local         |

  Scenario: all configured in config file
    And the configuration file:
      """
      [branches]
      main = "main"
      perennials = [ "public", "staging" ]
      perennial-regex = "^release-"
      feature-regex = "^user-.*$"
      contribution-regex = "^renovate/"
      observed-regex = "^dependabot/"
      unknown-type = "observed"

      [create]
      share-new-branches = "push"
      stash = false

      [hosting]
      forge-type = "github"
      github-connector = "gh"
      origin-hostname = "github.com"

      [ship]
      delete-tracking-branch = true
      strategy = "squash-merge"

      [sync]
      auto-resolve = false
      detached = true
      feature-strategy = "rebase"
      perennial-strategy = "ff-only"
      prototype-strategy = "compress"
      tags = false
      upstream = true
      """
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: contribution-1, contribution-2
        contribution regex: ^renovate/
        feature regex: ^user-.*$
        main branch: main
        observed branches: observed-1, observed-2
        observed regex: ^dependabot/
        parked branches: (none)
        perennial branches: public, staging
        perennial regex: ^release-
        prototype branches: (none)
        unknown branch type: observed

      Configuration:
        offline: no

      Create:
        new branch type: (not set)
        share new branches: push
        stash uncommitted changes: no

      Hosting:
        development remote: origin
        forge type: github
        origin hostname: github.com
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector type: gh
        GitHub token: (not set)
        GitLab connector type: (not set)
        GitLab token: (not set)

      Ship:
        delete tracking branch: yes
        ship strategy: squash-merge

      Sync:
        auto-resolve phantom conflicts: no
        run detached: yes
        run pre-push hook: yes
        feature sync strategy: rebase
        perennial sync strategy: ff-only
        prototype sync strategy: compress
        sync tags: no
        sync with upstream: yes
      """
