@smoke
Feature: require minimum Git version

  Scenario Outline: using an unsupported Git Version
    Given Git has version "2.6.2"
    When I run "git-town <COMMAND>"
    Then it prints the error:
      """
      this app requires Git 2.7.0 or higher
      """

    Examples:
      | COMMAND           |
      | aliases true      |
      | append foo        |
      | config            |
      | diff-parent       |
      | hack foo          |
      | kill              |
      | offline           |
      | propose           |
      | prepend foo       |
      | rename-branch foo |
      | repo              |
      | set-parent        |
      | ship              |
      | sync              |

  Scenario Outline: not requiring Git
    Given Git has version "2.6.2"
    When I run "git-town <COMMAND>"
    Then it runs no commands

    Examples:
      | COMMAND          |
      |                  |
      | completions bash |
      | help             |
      | --version        |
