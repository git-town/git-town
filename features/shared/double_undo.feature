Feature: no double undo

  Scenario:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given the current branch is "feature"
    And I run "git-town kill"
    And I run "git-town undo"
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
