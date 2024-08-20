Feature: append a new feature branch in a clean workspace using the "compress" sync strategy

  Background:
    Given a Git repo with origin
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE           |
      | existing | local, origin | existing commit 1 |
      | existing | local, origin | existing commit 2 |
    And the current branch is "existing"
    And Git Town setting "sync-feature-strategy" is "compress"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town append new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                                  |
      | existing | git fetch --prune --tags                 |
      |          | git checkout main                        |
      | main     | git rebase origin/main                   |
      |          | git checkout existing                    |
      | existing | git merge --no-edit --ff origin/existing |
      |          | git merge --no-edit --ff main            |
      |          | git reset --soft main                    |
      |          | git commit -m "existing commit 1"        |
      |          | git push --force-with-lease              |
      |          | git checkout -b new                      |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE           |
      | existing | local, origin | existing commit 1 |
      | new      | local         | existing commit 1 |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                                                   |
      | new      | git checkout existing                                     |
      | existing | git reset --hard {{ sha-before-run 'existing commit 2' }} |
      |          | git push --force-with-lease --force-if-includes           |
      |          | git branch -D new                                         |
    And the current branch is now "existing"
    And the initial commits exist
    And the initial lineage exists
