Feature: make another local feature branch a contribution branch

  Background:
    Given a local Git repo
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | local | feature | main   | local     |
    When I run "git-town contribute local"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "local" is local only - branches you want to contribute to must have a remote branch because they are per definition other people's branches
      """
    And branch "local" is still a feature branch
    And there are still no contribution branches
    And the current branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And branch "local" is still a feature branch
    And there are still no contribution branches
