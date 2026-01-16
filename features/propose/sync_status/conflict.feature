@skipWindows
Feature: merge conflict

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE       | FILE NAME        | FILE CONTENT   |
      | feature | local    | local commit  | conflicting_file | local content  |
      | feature | origin   | origin commit | conflicting_file | remote content |
    And the current branch is "feature"
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/feature |
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
    And no merge is now in progress
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
      |         | git push                                                           |
      |         | Finding proposal from feature into main ... none                   |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | feature | local, origin | local commit                                               |
      |         |               | origin commit                                              |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, commit, and continue
    Given I resolve the conflict in "conflicting_file"
    When I run "git commit --no-edit"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git push                                                           |
      |         | Finding proposal from feature into main ... none                   |
      |         | open https://github.com/git-town/git-town/compare/feature?expand=1 |
