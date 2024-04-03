Feature: enter a parent branch name when prompted

  Background:
    Given the branches "alpha" and "beta"
    And the current branch is "beta"

  Scenario: choose the default branch name
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                | KEYS  |
      | parent branch of beta | enter |
    Then this lineage exists now
      | BRANCH | PARENT |
      | beta   | main   |

  Scenario: choose other branches
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                 | KEYS       |
      | parent branch of beta  | down enter |
      | parent branch of alpha | enter      |
    And this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | beta   | alpha  |

  Scenario: choose "<none> (make a perennial branch)"
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                | KEYS     |
      | parent branch of beta | up enter |
    Then the perennial branches are now "beta"

  Scenario: enter the parent for several branches
    When I run "git-town sync --all" and enter into the dialog:
      | DIALOG                 | KEYS  |
      | parent branch of alpha | enter |
      | parent branch of beta  | enter |
    Then this lineage exists now
      | BRANCH | PARENT |
      | alpha  | main   |
      | beta   | main   |
