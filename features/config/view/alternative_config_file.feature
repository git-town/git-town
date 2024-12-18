Feature: show the configuration when using an alternative config file

  Background:
    Given a Git repo with origin

  Scenario: all configured in config file with alternative filename
    And file ".git-town.toml" with content
      """
      [branches]
      main = "main"
      perennials = [ "public", "staging" ]
      perennial-regex = "^release-"
      default-type = "observed"
      feature-regex = "^user-.*$"
      contribution-regex = "^renovate/"
      observed-regex = "^dependabot/"

      [create]
      push-new-branches = true

      [hosting]
      platform = "github"
      origin-hostname = "github.com"

      [ship]
      delete-tracking-branch = true
      strategy = "squash-merge"

      [sync]
      feature-strategy = "rebase"
      perennial-strategy = "merge"
      prototype-strategy = "compress"
      tags = false
      upstream = true
      """
    When I run "git-town config"
    Then Git Town prints:
      """
      Branches:
        contribution branches: (none)
        contribution regex: ^renovate/
        default branch type: observed
        feature regex: ^user-.*$
        main branch: main
        observed branches: (none)
        observed regex: ^dependabot/
        parked branches: (none)
        perennial branches: public, staging
        perennial regex: ^release-

      Configuration:
        offline: no

      Create:
        new branch type: feature
        push new branches: yes

      Hosting:
        hosting platform: github
        hostname: github.com
        GitHub token: (not set)
        GitLab token: (not set)
        Gitea token: (not set)

      Ship:
        delete the tracking branch: yes
        strategy: squash-merge

      Sync:
        run pre-push hook: yes
        sync-feature strategy: rebase
        sync-perennial strategy: merge
        sync-prototype strategy: compress
        sync tags: no
        sync with upstream: yes
      """
