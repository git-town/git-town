Feature: view changes made on the current feature branch

  Scenario: feature branch
    Given my repo has a feature branch "feature"
    And I am on the "feature" branch
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |

  Scenario: child branch
    Given my repo has a feature branch "feature-1"
    And my repo has a feature branch "feature-1A" as a child of "feature-1"
    And I am on the "feature-1A" branch
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH     | COMMAND                        |
      | feature-1A | git diff feature-1..feature-1A |
