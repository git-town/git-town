Feature: view changes made on another branch

  Background:
    Given my repo has a feature branch "feature-1"

  Scenario: feature branch
    When I run "git-town diff-parent feature-1"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git diff main..feature-1 |

  Scenario: child branch
    And my repo has a feature branch "feature-2" as a child of "feature-1"
    When I run "git-town diff-parent feature-2"
    Then it runs the commands
      | BRANCH | COMMAND                       |
      | main   | git diff feature-1..feature-2 |
