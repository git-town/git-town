Feature: undo changing offline mode

  Background:
    Given a Git repo with origin

  Scenario: undo enabling offline mode
    Given I ran "git-town offline on"
    And global Git Town setting "offline" is now "true"
    When I run "git-town undo"
    Then global Git Town setting "offline" now doesn't exist

  Scenario: undo disabling offline mode
    Given global Git Town setting "offline" is "true"
    And I ran "git-town offline off"
    And global Git Town setting "offline" is now "false"
    When I run "git-town undo"
    Then global Git Town setting "offline" is now "true"
