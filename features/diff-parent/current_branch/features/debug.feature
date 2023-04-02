Feature: display debug statistics

  Scenario: feature branch
    And the current branch is a feature branch "feature"
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |

  Scenario: child branch
    Given a feature branch "parent"
    And a feature branch "child" as a child of "parent"
    And the current branch is "child"
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH | COMMAND                |
      | child  | git diff parent..child |
