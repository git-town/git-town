Feature: ship the current feature branch with commit message in file

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the committed file "body.txt":
      """
      Commit message in file
      """
    And the current branch is "feature"
    When I run "git-town ship --message-file body.txt"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                |
      | feature | git fetch --prune --tags               |
      |         | git checkout main                      |
      | main    | git merge --squash --ff feature        |
      |         | git commit -m "Commit message in file" |
      |         | git push                               |
      |         | git push origin :feature               |
      |         | git branch -D feature                  |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | Commit message in file |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git revert {{ sha 'Commit message in file' }} |
      |        | git push                                      |
      |        | git branch feature {{ sha 'persisted file' }} |
      |        | git push -u origin feature                    |
      |        | git checkout feature                          |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                         |
      | main    | local, origin | Commit message in file          |
      |         |               | Revert "Commit message in file" |
      | feature | local, origin | feature commit                  |
      |         |               | persisted file                  |
