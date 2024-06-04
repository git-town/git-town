@skipWindows
Feature: merge conflict

  Background:
    Given the current branch is a local feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        | FILE NAME        | FILE CONTENT    |
      | main    | local, origin | main commit    | conflicting_file | main content    |
      | feature | local         | feature commit | conflicting_file | feature content |
    And tool "open" is installed
    And the origin is "git@github.com:git-town/git-town.git"
    When I run "git-town propose"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git fetch --prune --tags      |
      |         | git checkout main             |
      | main    | git rebase origin/main        |
      |         | git checkout feature          |
      | feature | git merge --no-edit --ff main |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      """
    And the current branch is still "feature"
    And a merge is now in progress

  @debug @this
  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND           |
      | feature | git merge --abort |
    And the current branch is still "feature"
    And no merge is in progress
    And the initial commits exist

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "feature"
    And a merge is now in progress

  Scenario: resolve and continue
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git commit --no-edit                                               |
      |         | git push -u origin feature                                         |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | main commit                      |
      | feature | local, origin | feature commit                   |
      |         |               | main commit                      |
      |         |               | Merge branch 'main' into feature |
    And these committed files exist now
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |

  Scenario: resolve, commit, and continue
    Given I resolve the conflict in "conflicting_file"
    When I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                                                            |
      | feature | git push -u origin feature                                         |
      | <none>  | open https://github.com/git-town/git-town/compare/feature?expand=1 |
    And "open" launches a new proposal with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """
    And the current branch is still "feature"
