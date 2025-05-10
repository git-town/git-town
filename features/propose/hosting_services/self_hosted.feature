@skipWindows
Feature: self-hosted service

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And a proposal for this branch does not exist

  Scenario Outline: self hosted
    Given tool "open" is installed
    And the origin is "git@self-hosted:git-town/git-town.git"
    And Git setting "git-town.forge-type" is "<PLATFORM>"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      | (none)  | Looking for proposal online ... ok |
      |         | open <PROPOSAL_URL>                |

    Examples:
      | PLATFORM  | PROPOSAL_URL                                                                                                                              |
      | bitbucket | https://self-hosted/git-town/git-town/pull-requests/new?source=feature&dest=git-town%2Fgit-town%3Amain                                    |
      | github    | https://self-hosted/git-town/git-town/compare/feature?expand=1                                                                            |
      | gitlab    | https://self-hosted/git-town/git-town/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |
# uncomment to test (creates online connection)
# | gitea     | https://self-hosted/git-town/git-town/compare/main...feature                                                                              |

  Scenario: GitLab with custom port
    Given the origin is "ssh://git@git.example.com:4022/a/b.git"
    And Git setting "git-town.forge-type" is "gitlab"
    And tool "open" is installed
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                                                              |
      | feature | git fetch --prune --tags                                                                                                             |
      | (none)  | Looking for proposal online ... ok                                                                                                   |
      |         | open https://git.example.com/a/b/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |
