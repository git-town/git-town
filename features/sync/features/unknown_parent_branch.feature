@messyoutput
Feature: enter a parent branch name when prompted

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE   | LOCATIONS     |
      | alpha | (none) | local, origin |
      | beta  | (none) | local, origin |
    And the current branch is "beta"

  Scenario: choose the default branch name
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                   | KEYS  |
      | parent branch for "beta" | enter |
    Then this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |

  Scenario: choose other branches
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                    | KEYS       |
      | parent branch for "beta"  | down enter |
      | parent branch for "alpha" | enter      |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | beta   | alpha  |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                   | KEYS     |
      | parent branch for "beta" | up enter |
    Then the perennial branches are now "beta"

  Scenario: enter the parent for several branches
    When I run "git-town sync --all" and enter into the dialog:
      | DIALOG                    | KEYS  |
      | parent branch for "beta"  | enter |
      | parent branch for "alpha" | enter |
    Then this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | beta   | main   |
