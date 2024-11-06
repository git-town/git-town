Feature: does not create prototyping branches this way

  Background:
    Given a Git repo with origin
    When I run "git-town prototype zonk"

  Scenario: result
    Then Git Town runs no commands
    And it prints the error:
      """
      there is no branch "zonk"
      """
    And there are still no prototype branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And it prints:
      """
      nothing to undo
      """
