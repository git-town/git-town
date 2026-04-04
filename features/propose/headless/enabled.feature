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

  Scenario: --headless flag with no existing proposal
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

  Scenario: --headless flag with existing proposal
    Given the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | URL                                           |
      | 1  | feature       | main          | https://github.com/git-town/git-town/pull/123 |
    When I run "git-town propose --headless"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                        |
      | feature | git fetch --prune --tags                                                       |
      |         | Finding proposal from feature into main ... #1 (Proposal from feature to main) |
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/pull/123
      """
    And the initial proposals exist now

  Scenario: headless configured via Git setting with no existing proposal
    Given Git setting "git-town.headless" is "true"
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

  Scenario: headless configured via Git setting with existing proposal
    Given the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | URL                                           |
      | 1  | feature       | main          | https://github.com/git-town/git-town/pull/123 |
    And Git setting "git-town.headless" is "true"
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                        |
      | feature | git fetch --prune --tags                                                       |
      |         | Finding proposal from feature into main ... #1 (Proposal from feature to main) |
    And Git Town prints:
      """
      Please open in a browser: https://github.com/git-town/git-town/pull/123
      """
    And the initial proposals exist now
