Feature: ship-delete-tracking-branch disabled when using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-delete-tracking-branch" is "false"
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "feature"
    When I run "git-town ship" and close the editor
    And origin deletes the "feature" branch

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                             |
      | feature | git fetch --prune --tags            |
      |         | git checkout main                   |
      | main    | git merge --no-ff --edit -- feature |
      |         | git push                            |
      |         | git branch -D feature               |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no lineage exists now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | feature commit         |
      |        |               | Merge branch 'feature' |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
      |        | git checkout feature                          |
    And the branches are now
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | origin     | main          |
    And the initial lineage exists now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | feature commit         |
      |        |               | Merge branch 'feature' |
