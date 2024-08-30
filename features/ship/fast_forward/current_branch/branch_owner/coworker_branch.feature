Feature: ship a coworker's feature branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE         | AUTHOR                          |
      | feature | local, origin | coworker commit | coworker <coworker@example.com> |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git fetch --prune --tags    |
      |         | git checkout main           |
      | main    | git merge --ff-only feature |
      |         | git push                    |
      |         | git push origin :feature    |
      |         | git branch -D feature       |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE         | AUTHOR                          |
      | main   | local, origin | coworker commit | coworker <coworker@example.com> |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                        |
      | main   | git branch feature {{ sha 'coworker commit' }} |
      |        | git push -u origin feature                     |
      |        | git checkout feature                           |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         |
      | main    | local, origin | coworker commit |
      | feature | local, origin | coworker commit |
    And the initial branches and lineage exist
