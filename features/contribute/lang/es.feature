Feature: make the current feature branch a contribution branch in Spanish

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town contribute" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      branch feature is now a contribution branch
      """
    And branch "feature" now has type "contribution"

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs no commands
    And branch "feature" now has type "feature"
