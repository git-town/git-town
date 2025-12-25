Feature: ship a branch that has the same name as a folder in the codebase

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE        | FILE NAME     |
      | main   | local, origin | feature commit | test/file.txt |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | test | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | test   | local, origin | commit 1 | shippable | content      |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "test"
    When I run "git-town ship -m testing"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                      |
      | test   | git fetch --prune --tags     |
      |        | git checkout main            |
      | main   | git merge --squash --ff test |
      |        | git commit -m testing        |
      |        | git push                     |
      |        | git push origin :test        |
      |        | git branch -D test           |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
      |        |               | testing        |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                              |
      | main   | git revert {{ sha 'testing' }}       |
      |        | git push                             |
      |        | git branch test {{ sha 'commit 1' }} |
      |        | git push -u origin test              |
      |        | git checkout test                    |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE          |
      | main   | local, origin | feature commit   |
      |        |               | testing          |
      |        |               | Revert "testing" |
      | test   | local, origin | commit 1         |
