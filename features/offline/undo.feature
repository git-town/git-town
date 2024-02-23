Feature: undo changing offline mode

  @this
  Scenario: undo enabling offline mode
    Given I run "git-town offline on"
    Then global Git Town setting "offline" is now "true"
    When I run "git-town undo"
    Then global Git Town setting "offline" now doesn't exist
