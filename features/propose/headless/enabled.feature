@skipWindows
Feature: headless mode

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And tool "open" is installed

  Scenario: --headless flag
    When I run "git-town propose --headless"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                          |
      | feature | git fetch --prune --tags                         |
      |         | Finding proposal from feature into main ... none |
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And the initial branches and lineage exist now

  Scenario: propose-headless configured via Git setting
    Given Git setting "git-town.propose-headless" is "true"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                          |
      | feature | git fetch --prune --tags                         |
      |         | Finding proposal from feature into main ... none |
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And the initial branches and lineage exist now
