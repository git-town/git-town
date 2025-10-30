Feature: sync the current feature branch without a tracking branch using the "compress" strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                |
      | main    | local    | local main commit      |
      |         | origin   | origin main commit     |
      | feature | local    | local feature commit 1 |
      | feature | local    | local feature commit 2 |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git push                                          |
      |         | git checkout feature                              |
      | feature | git merge --no-edit --ff main                     |
      |         | git reset --soft main --                          |
      |         | git commit -m "local feature commit 1"            |
      |         | git push -u origin feature                        |
    And the branches are now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                |
      | main    | local, origin | origin main commit     |
      |         |               | local main commit      |
      | feature | local, origin | local feature commit 1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                     |
      | feature | git reset --hard {{ sha-initial 'local feature commit 2' }} |
      |         | git push origin :feature                                    |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                |
      | main    | local, origin | origin main commit     |
      |         |               | local main commit      |
      | feature | local         | local feature commit 1 |
      |         |               | local feature commit 2 |
