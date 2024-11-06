Feature: sync the current feature branch (in a local repo)

  Background:
    Given a local Git repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | main    | local    | main commit    |
      | feature | local    | feature commit |
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git merge --no-edit --ff main |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | main commit                      |
      | feature | local    | feature commit                   |
      |         |          | Merge branch 'main' into feature |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                     |
      | feature | git reset --hard {{ sha 'feature commit' }} |
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial branches and lineage exist now
