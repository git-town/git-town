@skipWindows
Feature: set a custom browser via the CLI

  Background:
    Given a Git repo with origin
    And the origin is "https://github.com/git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And tool "firefox" is installed
    And tool "open" is installed
    When I run "git-town propose --browser=firefox"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                               |
      | feature | git fetch --prune --tags                                              |
      |         | Finding proposal from feature into main ... none                      |
      |         | firefox https://github.com/git-town/git-town/compare/feature?expand=1 |
