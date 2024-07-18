Feature: does not diff non-existing branch

  Scenario: non-existing branch
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town diff-parent non-existing"
    Then it runs no commands
    And it prints the error:
      """
      there is no branch "non-existing"
      """
