Feature: prototype another prototype branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And an uncommitted file
    And the current branch is "prototype"
    When I run "git-town prototype prototype"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      branch "prototype" is already a prototype branch
      """
    And the prototype branches are now "prototype"
    And the current branch is still "prototype"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the prototype branches are now "prototype"
    And the current branch is still "prototype"
    And the uncommitted file still exists
    And the initial branches and lineage exist
