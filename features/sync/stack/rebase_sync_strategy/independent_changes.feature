Feature: stacked changes where each branch makes independent changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                        |
      | child  | git fetch --prune --tags                                                       |
      |        | git checkout main                                                              |
      | main   | git -c rebase.updateRefs=false rebase origin/main                              |
      |        | git push                                                                       |
      |        | git checkout parent                                                            |
      | parent | git push --force-with-lease --force-if-includes                                |
      |        | git -c rebase.updateRefs=false rebase origin/parent                            |
      |        | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }}   |
      |        | git push --force-with-lease --force-if-includes                                |
      |        | git checkout child                                                             |
      | child  | git push --force-with-lease --force-if-includes                                |
      |        | git -c rebase.updateRefs=false rebase origin/child                             |
      |        | git -c rebase.updateRefs=false rebase --onto parent {{ sha 'initial commit' }} |
      |        | git push --force-with-lease --force-if-includes                                |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | parent | local, origin | origin parent commit |
      |        |               | local parent commit  |
      | child  | local, origin | origin child commit  |
      |        |               | local child commit   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                      |
      | child  | git reset --hard {{ sha-initial 'local child commit' }}                                      |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin child commit' }}:child   |
      |        | git checkout parent                                                                          |
      | parent | git reset --hard {{ sha-initial 'local parent commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin parent commit' }}:parent |
      |        | git checkout child                                                                           |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | parent | local         | local parent commit  |
      |        | origin        | origin parent commit |
      | child  | local         | local child commit   |
      |        | origin        | origin child commit  |
