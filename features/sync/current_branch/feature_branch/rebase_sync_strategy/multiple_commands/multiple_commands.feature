Feature: running a sync after running another Git Town command

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
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
    And the current branch is "child"
    And I ran "git-town sync"
    And I ran "git-town hack new"
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                             |
      | child  | git fetch --prune --tags                                                            |
      |        | git checkout parent                                                                 |
      | parent | git -c rebase.updateRefs=false rebase --onto main {{ sha 'local main commit' }}     |
      |        | git checkout child                                                                  |
      | child  | git -c rebase.updateRefs=false rebase --onto parent {{ sha 'local parent commit' }} |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | child  | local, origin | origin child commit  |
      |        |               | local child commit   |
      | parent | local, origin | origin parent commit |
      |        |               | local parent commit  |
    And these branches exist now
      | REPOSITORY | BRANCHES                 |
      | local      | main, child, new, parent |
      | origin     | main, child, parent      |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | child  | local, origin | origin child commit  |
      |        |               | local child commit   |
      | parent | local, origin | origin parent commit |
      |        |               | local parent commit  |
    And these branches exist now
      | REPOSITORY | BRANCHES                 |
      | local      | main, child, new, parent |
      | origin     | main, child, parent      |
