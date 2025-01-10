Feature: syncing a stack that contains an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE   | LOCATIONS |
      | observed | (none) | origin    |
    And the current branch is "main"
    And I ran "git-town observe observed"
    And the commits
      | BRANCH   | LOCATION | MESSAGE    |
      | observed | origin   | new commit |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    # And inspect the repo
    When I run "git-town sync --stack"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git push --tags          |
    And all branches are now synchronized
    And the current branch is still "main"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | child  | local, origin | origin child commit  |
      |        |               | origin parent commit |
      |        |               | origin main commit   |
      |        |               | local main commit    |
      |        |               | local parent commit  |
      |        |               | local child commit   |
      | parent | local, origin | origin parent commit |
      |        |               | origin main commit   |
      |        |               | local main commit    |
      |        |               | local parent commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                         |
      | child  | git reset --hard {{ sha-before-run 'local child commit' }}                                      |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin child commit' }}:child   |
      |        | git checkout parent                                                                             |
      | parent | git reset --hard {{ sha-before-run 'local parent commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin parent commit' }}:parent |
      |        | git checkout child                                                                              |
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | child  | local         | local child commit   |
      |        | origin        | origin child commit  |
      | parent | local         | local parent commit  |
      |        | origin        | origin parent commit |
    And the initial branches and lineage exist now
