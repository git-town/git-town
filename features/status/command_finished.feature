Feature:

  @this
  Scenario: Git Town command ran successfully
    Given I run "git-town sync"
    When I run "git-town status"
    Then it prints:
      """
      The previous Git Town command (sync) finished successfully.
      You can run "git town undo" to undo it.
      """
