Feature: append with configured branch-prefix via Git metadata

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And Git setting "git-town.branch-prefix" is "kg-"
    And the current branch is "feature-1"

  Scenario Outline:
    When I run "git-town append <BRANCH_NAME>"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                      |
      | feature-1 | git fetch --prune --tags     |
      |           | git checkout -b kg-feature-2 |
    And the current branch is now "kg-feature-2"
    And this lineage exists now
      """
      main
        feature-1
          kg-feature-2
      """

    Examples:
      | BRANCH_NAME  |
      | feature-2    |
      | kg-feature-2 |

  Scenario: undo
    Given I ran "git-town append feature-2"
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                    |
      | kg-feature-2 | git checkout feature-1     |
      | feature-1    | git branch -D kg-feature-2 |
    And the current branch is now "feature-1"
    And this lineage exists now
      """
      main
        feature-1
      """
