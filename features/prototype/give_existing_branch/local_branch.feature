Feature: prototype another local feature branch

  Background:
    Given a local Git repo
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    When I run "git-town prototype feature"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
      branch "feature" is now a prototype branch
      """
    And the prototype branches are now "feature"

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And branch "feature" is now a feature branch
    And there are now no prototype branches
