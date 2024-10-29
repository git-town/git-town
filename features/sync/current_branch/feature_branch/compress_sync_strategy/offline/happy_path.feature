Feature: sync the current feature branch using the "compress" strategy in offline mode

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git Town setting "sync-feature-strategy" is "compress"
    And offline mode is enabled
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE                |
      | main    | local    | local main commit      |
      |         | origin   | origin main commit     |
      | feature | local    | local feature commit 1 |
      |         | local    | local feature commit 2 |
      |         | origin   | origin feature commit  |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git reset --soft main                   |
      |         | git commit -m "local feature commit 1"  |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE                |
      | main    | local    | local main commit      |
      |         | origin   | origin main commit     |
      | feature | local    | local feature commit 1 |
      |         | origin   | origin feature commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                        |
      | feature | git reset --hard {{ sha-before-run 'local feature commit 2' }} |
    And the current branch is still "feature"
    And the initial commits exist now
    And the initial branches and lineage exist now
