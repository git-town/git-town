@skipWindows
Feature: propose with embedded lineage

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.proposals-show-lineage" is "cli"
    And the current branch is "feature"
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                                                                                                                                                                                                                                                                                                                                                      |
      | feature | git fetch --prune --tags                                                                                                                                                                                                                                                                                                                                                                                                     |
      |         | Finding proposal from feature into main ... none                                                                                                                                                                                                                                                                                                                                                                             |
      |         | Finding proposal from feature into main ... none                                                                                                                                                                                                                                                                                                                                                                             |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1&body=%3C%21--+branch-stack-start+--%3E%0A%0A-------------------------%0A-+main%0A++-+feature+%3Apoint_left%3A%0A%0A%3Csup%3E%5BStack%5D%28https%3A%2F%2Fwww.git-town.com%2Fhow-to%2Fgithub-actions-breadcrumb.html%29+generated+by+%5BGit+Town%5D%28https%3A%2F%2Fgithub.com%2Fgit-town%2Fgit-town%29%3C%2Fsup%3E%0A%0A%3C%21--+branch-stack-end+--%3E%0A |
      |         | Finding all proposals for feature ... none                                                                                                                                                                                                                                                                                                                                                                                   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial branches exist now
