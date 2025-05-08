Feature: ship an omni-branch via the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "always-merge"
    And local Git setting "color.ui" is "always"
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
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | feature commit         |
      |        |               | Merge branch 'feature' |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                  |
      | main   | git branch feature {{ sha-before-run 'feature commit' }} |
      |        | git push -u origin feature                               |
      |        | git checkout feature                                     |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                |
      | main   | local, origin | feature commit         |
      |        |               | Merge branch 'feature' |
    And the initial branches and lineage exist now
