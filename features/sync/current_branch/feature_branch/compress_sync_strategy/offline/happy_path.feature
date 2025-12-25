Feature: sync the current feature branch using the "compress" strategy in offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                |
      | main    | local    | local main commit      |
      |         | origin   | origin main commit     |
      | feature | local    | local feature commit 1 |
      |         | local    | local feature commit 2 |
      |         | origin   | origin feature commit  |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And offline mode is enabled
    And the current branch is "feature"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git checkout feature                              |
      | feature | git merge --no-edit --ff main                     |
      |         | git merge --no-edit --ff origin/feature           |
      |         | git reset --soft main --                          |
      |         | git commit -m "local feature commit 1"            |
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                |
      | main    | local    | local main commit      |
      |         | origin   | origin main commit     |
      | feature | local    | local feature commit 1 |
      |         | origin   | origin feature commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                     |
      | feature | git reset --hard {{ sha-initial 'local feature commit 2' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
