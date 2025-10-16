Feature: ship the supplied feature branch from a subfolder using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
      | other   | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And Git setting "git-town.ship-strategy" is "always-merge"
    And the current branch is "other"
    And a folder "new_folder"
    When I run "git-town ship -m 'feature done' feature" in the "new_folder" folder

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                        |
      | other  | git fetch --prune --tags                       |
      |        | git checkout main                              |
      | main   | git merge --no-ff -m "feature done" -- feature |
      |        | git push                                       |
      |        | git checkout other                             |
      | other  | git branch -D feature                          |
    And this lineage exists now
      """
      main
        other
      """
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, other |
      | origin     | main        |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
      |        |               | feature done   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | other  | git branch feature {{ sha 'feature commit' }} |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE        |
      | main   | local, origin | feature commit |
      |        |               | feature done   |
