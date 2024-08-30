Feature: observe the current local feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town observe"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "feature" is local only - branches you want to observe must have a remote branch because they are per definition other people's branches
      """
    And the current branch is still "feature"
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And there are still no observed branches
