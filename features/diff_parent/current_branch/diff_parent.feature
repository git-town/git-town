@smoke
Feature: view changes made on the current feature branch

  Background:
    Given a Git repo clone

  Scenario: feature branch
    Given the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH  | COMMAND                |
      | feature | git diff main..feature |

  Scenario: child branch
    Given the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | parent | feature | main   | local     |
      | child  | feature | parent | local     |
    And the current branch is "child"
    When I run "git-town diff-parent"
    Then it runs the commands
      | BRANCH | COMMAND                |
      | child  | git diff parent..child |
