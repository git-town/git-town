Feature: append a new branch when feature branches are configured

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | existing | local, origin | existing commit |
    And the current branch is "existing"
    And Git setting "git-town.new-branch-type" is "feature"
    And Git setting "git-town.default-branch-type" is "contribution"
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout -b new main |
    And branch "new" now has type "feature"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | new    |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the initial commits exist now
    And the initial lineage exists now
