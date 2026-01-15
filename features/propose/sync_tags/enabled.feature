Feature: don't sync tags while proposing

  Background:
    Given a Git repo with origin
    And the origin is "ssh://git@github.com/git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git setting "git-town.sync-tags" is "false"
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "feature"
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git fetch --prune --no-tags                                        |
      |         | Finding proposal from feature into main ... ok                     |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial commits exist now
    And the initial tags exist now
