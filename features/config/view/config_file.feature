Feature: display information from config file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE         | PARENT | LOCATIONS     |
      | contribution-1 | contribution |        | local, origin |
      | contribution-2 | contribution |        | local, origin |
      | observed-1     | observed     |        | local, origin |
      | observed-2     | observed     |        | local, origin |
      | parked-1       | parked       | main   | local         |
      | parked-2       | parked       | main   | local         |
      | perennial-1    | perennial    |        | local         |
      | perennial-2    | perennial    |        | local         |
      | prototype-1    | prototype    | main   | local         |
      | prototype-2    | prototype    | main   | local         |

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
      order = "desc"
      display-types = "all"

      [create]
      branch-prefix = "acme-"
      share-new-branches = "push"
      stash = false

      [hosting]
      browser = "chrome"
      forge-type = "github"
      github-connector = "gh"
      origin-hostname = "github.com"

      [propose]
      breadcrumb = "stacks"

      [ship]
      delete-tracking-branch = true
      ignore-uncommitted = true
      strategy = "squash-merge"

      [sync]
      auto-resolve = false
      auto-sync = false
      detached = true
      feature-strategy = "rebase"
      perennial-strategy = "ff-only"
      prototype-strategy = "compress"
      push-branches = false
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
        parked branches: parked-1, parked-2
        perennial branches: public, staging
        perennial regex: ^release-
        prototype branches: prototype-1, prototype-2
        unknown branch type: observed
        order: desc
        display types: all branch types

      Configuration:
        offline: no
        git user name: user
        git user email: email@example.com

      Create:
        branch prefix: acme-
        new branch type: (not set)
        share new branches: push
        stash uncommitted changes: no

      Hosting:
        browser: chrome
        development remote: origin
        forge type: github
        origin hostname: github.com
        Bitbucket username: (not set)
        Bitbucket app password: (not set)
        Forgejo token: (not set)
        Gitea token: (not set)
        GitHub connector: gh
        GitHub token: (not set)
        GitLab connector: (not set)
        GitLab token: (not set)

      Propose:
        breadcrumb: stacks
        breadcrumb direction: down

      Ship:
        delete tracking branch: yes
        ignore uncommitted changes: yes
        ship strategy: squash-merge

      Sync:
        auto-resolve phantom conflicts: no
        auto-sync: no
        run detached: yes
        run pre-push hook: yes
        feature sync strategy: rebase
        perennial sync strategy: ff-only
        prototype sync strategy: compress
        push branches: no
        sync tags: no
        sync with upstream: yes
        auto-resolve phantom conflicts: no
      """
