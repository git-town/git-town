Feature: view changes made on the current feature branch

  Scenario: feature branch
    Given my repo has a feature branch "feature"
    And I am on the "feature" branch
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |

  Scenario: child branch
    Given my repo has a feature branch "parent"
    And my repo has a feature branch "child" as a child of "parent"
    And I am on the "child" branch
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH | COMMAND                |
      | child  | git diff parent..child |
