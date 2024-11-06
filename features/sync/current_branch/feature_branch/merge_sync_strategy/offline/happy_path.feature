Feature: offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And offline mode is enabled
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | local main commit                |
      |         | origin   | origin main commit               |
      | feature | local    | local feature commit             |
      |         |          | Merge branch 'main' into feature |
      |         | origin   | origin feature commit            |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git reset --hard {{ sha 'local feature commit' }} |
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial branches and lineage exist now
