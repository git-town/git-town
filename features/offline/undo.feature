Feature: undo changing offline mode

  Scenario: undo enabling offline mode
    Given I run "git-town offline on"
    Then global Git Town setting "offline" is now "true"
    When I run "git-town undo"
    Then global Git Town setting "offline" now doesn't exist

  Scenario: undo disabling offline mode
    Given global Git Town setting "offline" is "true"
    And I run "git-town offline off"
    Then global Git Town setting "offline" is now "false"
    When I run "git-town undo"
    Then global Git Town setting "offline" is now "true"
