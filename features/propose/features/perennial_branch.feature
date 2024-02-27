Feature: Cannot create proposals for perennial branches

  Background:
    Given the current branch is a perennial branch "perennial"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND |
    And it prints the error:
      """
      xxx
      """
