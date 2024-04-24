Feature: sync the current feature branch (in a local repo)

  Background:
    Given my repo does not have an origin
    And the current branch is a local feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | main    | local    | main commit    |
      | feature | local    | feature commit |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git merge --no-edit --ff main |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | main commit                      |
      | feature | local    | feature commit                   |
      |         |          | main commit                      |
      |         |          | Merge branch 'main' into feature |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                     |
      | feature | git reset --hard {{ sha 'feature commit' }} |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
