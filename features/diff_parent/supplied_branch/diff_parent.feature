@smoke
Feature: view changes made on another branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | alpha | feature | main   | local     |

  Scenario: feature branch
    When I run "git-town diff-parent alpha"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | main   | git diff main..alpha |

  Scenario: child branch
    Given the branches
      | NAME | TYPE    | PARENT | LOCATIONS |
      | beta | feature | alpha  | local     |
    When I run "git-town diff-parent beta"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | main   | git diff alpha..beta |
