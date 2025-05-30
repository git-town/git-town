@smoke
Feature: require minimum Git version

  Background:
    Given a Git repo with origin
    And Git has version "2.29.2"

  Scenario Outline: using an unsupported Git Version
    When I run "git-town <COMMAND>"
    Then Git Town prints the error:
      """
      this app requires Git 2.30 or higher
      """

    Examples:
      | COMMAND      |
      | append foo   |
      | config       |
      | config setup |
      | diff-parent  |
      | hack foo     |
      | delete       |
      | offline      |
      | propose      |
      | prepend foo  |
      | rename foo   |
      | repo         |
      | set-parent   |
      | ship         |
      | sync         |

  Scenario Outline: not requiring Git
    When I run "git-town <COMMAND>"
    Then Git Town runs no commands

    Examples:
      | COMMAND          |
      |                  |
      | completions bash |
      | help             |
      | --version        |
