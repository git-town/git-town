Feature: sync a branch that contains a file with the same name

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           |
      | main    | local    | local main commit     | local main file     |
      |         | origin   | origin main commit    | origin main file    |
      | feature | local    | local feature commit  | feature             |
      |         | origin   | origin feature commit | origin feature file |
    And the current branch is "feature"
    And inspect the repo
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git push                                          |
      |         | git checkout feature                              |
    And Git Town prints the error:
      """
      git log --no-merges --format=%H main ^feature
      """
    And Git Town prints the error:
      """
      ambiguous argument 'feature': both revision and filename
      """
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, origin | origin main commit                                         |
      |         |               | local main commit                                          |
      | feature | local, origin | local feature commit                                       |
      |         |               | Merge branch 'main' into feature                           |
      |         |               | origin feature commit                                      |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | origin main commit    |
      |         |               | local main commit     |
      | feature | local         | local feature commit  |
      |         | origin        | origin feature commit |
    And the initial branches and lineage exist now
