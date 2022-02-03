Feature: view changes made on another branch

  Scenario: feature branch
    Given my repo has a feature branch "feature-1"
    When I run "git-town diff-parent feature-1"
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git diff main..feature-1 |

  Scenario: child branch
    Given my repo has a feature branch "feature-1"
    And my repo has a feature branch "feature-2" as a child of "feature-1"
    When I run "git-town diff-parent feature-2"
    Then it runs the commands
      | BRANCH | COMMAND                       |
      | main   | git diff feature-1..feature-2 |
