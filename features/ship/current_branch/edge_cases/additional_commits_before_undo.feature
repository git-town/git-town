Feature: can undo a ship even after additional commits to the main branch

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    When I run "git-town ship -m done"
    And I add commit "additional commit" to the "main" branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git revert {{ sha 'done' }}                   |
      |        | git push                                      |
      |        | git branch feature {{ sha 'feature commit' }} |
      |        | git push -u origin feature                    |
      |        | git checkout feature                          |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE           |
      | main    | local, origin | done              |
      |         |               | additional commit |
      |         |               | Revert "done"     |
      | feature | local, origin | feature commit    |
    And the initial branches and lineage exist
