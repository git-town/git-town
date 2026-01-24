Feature: propose single stack, show single stacks disabled

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  |
      | branch-1 | local, origin | commit 1 |
    And Git setting "git-town.proposal-breadcrumb" is "cli"
    And Git setting "git-town.proposal-breadcrumb-single" is "no"
    And the current branch is "branch-1"
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                             |
      | branch-1 | git fetch --prune --tags                                            |
      |          | Finding proposal from branch-1 into main ... none                   |
      |          | open https://github.com/git-town/git-town/compare/branch-1?expand=1 |
      |          | Finding all proposals for branch-1 ... none                         |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial branches exist now
