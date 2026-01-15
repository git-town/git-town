@skipWindows
Feature: use a SSH identity

  Scenario Outline: ssh identity
    Given a Git repo with origin
    And the origin is "git@my-ssh-identity:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git setting "git-town.hosting-origin-hostname" is "<ORIGIN_HOSTNAME>"
    And the current branch is "feature"
    And tool "open" is installed
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                        |
      | feature | git fetch --prune --tags                       |
      |         | Finding proposal from feature into main ... ok |
      |         | open <PROPOSAL_URL>                            |

    Examples:
      | ORIGIN_HOSTNAME | PROPOSAL_URL                                                                                                                             |
      | bitbucket.org   | https://bitbucket.org/git-town/git-town/pull-requests/new?source=feature&dest=git-town%2Fgit-town%3Amain                                 |
      | github.com      | https://github.com/git-town/git-town/compare/feature?expand=1                                                                            |
      | gitlab.com      | https://gitlab.com/git-town/git-town/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main |
