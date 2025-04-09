@skipWindows
Feature: merge conflict

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | main commit    | conflicting_file | main content    |
      | feature | local         | feature commit | conflicting_file | feature content |
    And the current branch is "feature"
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exists
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      | (none)  | Looking for proposal online ... ok |
      | feature | git merge --no-edit --ff main      |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
    And no merge is in progress
    And the initial commits exist now

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a merge is now in progress

  Scenario: resolve and continue
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git commit --no-edit                                               |
      |         | git push -u origin feature                                         |
      | (none)  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | main commit                      |
      | feature | local, origin | feature commit                   |
      |         |               | Merge branch 'main' into feature |
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, commit, and continue
    Given I resolve the conflict in "conflicting_file"
    When I run "git commit --no-edit"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git push -u origin feature                                         |
      | (none)  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
