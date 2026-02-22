@skipWindows
Feature: init without TTY

  Background:
    Given a Git repo with origin
    And Git Town is not configured
    And local Git setting "init.defaultbranch" is "main"
    When I run "git-town init" in a non-TTY shell

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      no interactive terminal available
      """
