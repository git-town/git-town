Feature: detached sync the current feature branch using the "rebase" feature sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    When I run "git-town sync --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                              |
      | feature | git fetch --prune --tags                             |
      |         | git push --force-with-lease --force-if-includes      |
      |         | git -c rebase.updateRefs=false rebase origin/feature |
      |         | git push --force-with-lease --force-if-includes      |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local         | local main commit     |
      |         | origin        | origin main commit    |
      | feature | local, origin | origin feature commit |
      |         |               | local feature commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                        |
      | feature | git reset --hard {{ sha-initial 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin-initial 'origin feature commit' }}:feature |
    And the initial branches and lineage exist now
    And the initial commits exist now
