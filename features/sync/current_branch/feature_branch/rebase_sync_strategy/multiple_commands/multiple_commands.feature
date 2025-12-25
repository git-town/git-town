Feature: running a sync after running another Git Town command

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the commits
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main commit |
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | parent | local    | parent commit |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | child | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | child  | local    | child commit |
    And I ran "git-town hack new"
    And the current branch is "child"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git checkout parent                             |
      | parent | git push --force-with-lease --force-if-includes |
      |        | git checkout child                              |
      | child  | git push --force-with-lease --force-if-includes |
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, child, new, parent |
      | origin     | main, child, parent      |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | main commit   |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                            |
      | child  | git push --force-with-lease origin {{ sha 'parent commit' }}:child |
      |        | git push --force-with-lease origin {{ sha 'main commit' }}:parent  |
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, child, new, parent |
      | origin     | main, child, parent      |
    And the initial commits exist now
