@smoke
Feature: view changes made on another branch

  Background:
    Given a feature branch "alpha"

  Scenario: feature branch
    When I run "git-town diff-parent alpha"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | main   | git diff main..alpha |

  Scenario: child branch
    And a feature branch "beta" as a child of "alpha"
    When I run "git-town diff-parent beta"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | main   | git diff alpha..beta |
