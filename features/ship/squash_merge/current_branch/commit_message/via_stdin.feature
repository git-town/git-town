Feature: ship the current feature branch with commit message via STDIN

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
    When I pipe the following text into "git-town ship -f -":
      """
      Commit message via STDIN
      """

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                  |
      | feature | git fetch --prune --tags                 |
      |         | git checkout main                        |
      | main    | git merge --squash --ff feature          |
      |         | git commit -m "Commit message via STDIN" |
      |         | git push                                 |
      |         | git push origin :feature                 |
      |         | git branch -D feature                    |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                  |
      | main   | local, origin | Commit message via STDIN |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | main   | git revert {{ sha 'Commit message via STDIN' }} |
      |        | git push                                        |
      |        | git branch feature {{ sha 'feature commit' }}   |
      |        | git push -u origin feature                      |
      |        | git checkout feature                            |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                           |
      | main    | local, origin | Commit message via STDIN          |
      |         |               | Revert "Commit message via STDIN" |
      | feature | local, origin | feature commit                    |
