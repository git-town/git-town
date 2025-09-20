@smoke
Feature: view changes made on another branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | alpha | feature | main   | local     |
    And the current branch is "main"

  Scenario: feature branch
    When I run "git-town diff-parent alpha"
    Then Git Town runs the commands
      | BRANCH | COMMAND                          |
      | main   | git diff --merge-base main alpha |

  Scenario: child branch
    Given the branches
      | NAME | TYPE    | PARENT | LOCATIONS |
      | beta | feature | alpha  | local     |
    And the current branch is "main"
    When I run "git-town diff-parent beta"
    Then Git Town runs the commands
      | BRANCH | COMMAND                          |
      | main   | git diff --merge-base alpha beta |
