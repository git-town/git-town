Feature: ship a coworker's feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE         | AUTHOR                          |
      | feature | local, origin | coworker commit | coworker <coworker@example.com> |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship" and close the editor

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                             |
      | feature | git fetch --prune --tags            |
      |         | git checkout main                   |
      | main    | git merge --no-ff --edit -- feature |
      |         | git push                            |
      |         | git push origin :feature            |
      |         | git branch -D feature               |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                | AUTHOR                          |
      | main   | local, origin | coworker commit        | coworker <coworker@example.com> |
      |        |               | Merge branch 'feature' | user <email@example.com>        |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git branch feature {{ sha 'coworker commit' }} |
      |        | git push -u origin feature                     |
      |        | git checkout feature                           |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | coworker commit        |
      |        |               | Merge branch 'feature' |
    And the initial branches and lineage exist now
