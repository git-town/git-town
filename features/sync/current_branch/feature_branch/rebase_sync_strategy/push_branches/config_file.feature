Feature: disable pushing through the config file

  Background:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [branches]
      main = "main"

      [sync]
      push-branches = false
      feature-strategy = "rebase"
      """
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
    And Git Town is not configured
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                               |
      | child  | git fetch --prune --tags                                                              |
      |        | git checkout main                                                                     |
      | main   | git -c rebase.updateRefs=false rebase origin/main                                     |
      |        | git checkout parent                                                                   |
      | parent | git -c rebase.updateRefs=false rebase origin/parent                                   |
      |        | git -c rebase.updateRefs=false rebase --onto main {{ sha 'persisted config file' }}   |
      |        | git checkout child                                                                    |
      | child  | git -c rebase.updateRefs=false rebase origin/child                                    |
      |        | git -c rebase.updateRefs=false rebase --onto parent {{ sha 'persisted config file' }} |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        | local         | local main commit    |
      | parent | local         | origin parent commit |
      |        |               | local parent commit  |
      |        | origin        | origin parent commit |
      | child  | local         | origin child commit  |
      |        |               | local child commit   |
      |        | origin        | origin child commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                          |
      | child  | git reset --hard {{ sha 'local child commit' }}  |
      |        | git checkout main                                |
      | main   | git reset --hard {{ sha 'local main commit' }}   |
      |        | git checkout parent                              |
      | parent | git reset --hard {{ sha 'local parent commit' }} |
      |        | git checkout child                               |
    And the initial branches and lineage exist now
    And the initial commits exist now
