Feature: append a new feature branch in a dirty workspace using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           |
      | existing | local, origin | existing commit 1 |
      | existing | local, origin | existing commit 2 |
    And the current branch is "existing"
    And an uncommitted file
    And Git Town setting "sync-feature-strategy" is "compress"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND             |
      | existing | git add -A          |
      |          | git stash           |
      |          | git checkout -b new |
      | new      | git stash pop       |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE           |
      | existing | local, origin | existing commit 1 |
      |          |               | existing commit 2 |
      | new      | local         | existing commit 1 |
      |          |               | existing commit 2 |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git add -A            |
      |          | git stash             |
      |          | git checkout existing |
      | existing | git branch -D new     |
      |          | git stash pop         |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial lineage exists
