Feature: can undo a ship even after additional commits to the main branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "feature"
    When I run "git-town ship -m done"

  Scenario: add commit and undo
    When I add commit "additional commit" to the "main" branch
    And I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git revert {{ sha 'done' }}                   |
      |        | git push                                      |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
      |        | git checkout feature                          |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE           |
      | main    | local, origin | done              |
      |         |               | additional commit |
      |         |               | Revert "done"     |
      | feature | local, origin | feature commit    |
