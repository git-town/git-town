Feature: make a local branch a contribution branch

  Background:
    Given a local Git repo
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town contribute feature"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      Branches you want to contribute to must have a remote branch because they are per definition other people's branches.
      """
    And branch "feature" is still a feature branch
    And there are still no contribution branches
    And the current branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "main"
    And branch "feature" is still a feature branch
    And there are still no contribution branches
