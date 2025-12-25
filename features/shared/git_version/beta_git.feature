Feature: support Git compiled from development branch

  Background:
    Given a Git repo with origin
    And Git has version "2.47.GIT"

  Scenario Outline:
    When I run "git-town <COMMAND>"
    Then Git Town runs without errors

    Examples:
      | COMMAND    |
      | append foo |
      | config     |
      | hack foo   |
      | offline    |
      | sync       |
