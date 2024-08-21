Feature: cannot make the current local feature branch a contribution branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town contribute"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "feature" is local only - branches you want to contribute to must have a remote branch because they are per definition other people's branches
      """
    And the current branch is still "feature"
    And there are still no contribution branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And there are still no contribution branches
