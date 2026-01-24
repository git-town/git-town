Feature: propose single stack, show single stacks enabled

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
    And Git setting "git-town.proposals-show-lineage" is "cli"
    And Git setting "git-town.proposal-breadcrumb-single" is "yes"
    And the current branch is "branch-1"
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                                                                                                                                                                                                                                                                                                                                                        |
      | branch-1 | git fetch --prune --tags                                                                                                                                                                                                                                                                                                                                                                                                       |
      |          | Finding proposal from branch-1 into main ... none                                                                                                                                                                                                                                                                                                                                                                              |
      |          | open https://github.com/git-town/git-town/compare/branch-1?expand=1&body=%3C%21--+branch-stack-start+--%3E%0A%0A-------------------------%0A-+main%0A++-+branch-1+%3Apoint_left%3A%0A%0A%3Csup%3E%5BStack%5D%28https%3A%2F%2Fwww.git-town.com%2Fhow-to%2Fgithub-actions-breadcrumb.html%29+generated+by+%5BGit+Town%5D%28https%3A%2F%2Fgithub.com%2Fgit-town%2Fgit-town%29%3C%2Fsup%3E%0A%0A%3C%21--+branch-stack-end+--%3E%0A |
      |          | Finding all proposals for branch-1 ... none                                                                                                                                                                                                                                                                                                                                                                                    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial branches exist now
