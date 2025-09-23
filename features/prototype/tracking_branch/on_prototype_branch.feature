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
    And Git Town prints:
      """
      branch "migrate-task-20223" is already a prototype branch
      """
    And branch "migrate-task-20223" still has type "prototype"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "migrate-task-20223" still has type "prototype"
