Feature: prototype the current prototoype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME               | TYPE      | PARENT | LOCATIONS |
      | migrate-task-20223 | prototype | main   | local     |
    And the current branch is "migrate-task-20223"
    When I run "git-town prototype"

  Scenario: result
    Then Git Town runs no commands
    And it prints the error:
      """
      branch "migrate-task-20223" is already a prototype branch
      """
    And the prototype branches are still "migrate-task-20223"
    And the current branch is still "migrate-task-20223"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the prototype branches are still "migrate-task-20223"
    And the current branch is still "migrate-task-20223"
