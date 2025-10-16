Feature: hack a new branch while the main branch is active in another worktree

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE            |
      | main     | origin   | origin main commit |
      |          | local    | local main commit  |
      | existing | local    | existing commit    |
    And the current branch is "existing"
    And branch "main" is active in another worktree
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout -b new main |
    And this lineage exists now
      """
      main
        existing
        new
      """
    And these commits exist now
      | BRANCH   | LOCATION        | MESSAGE            |
      | main     | origin          | origin main commit |
      |          | worktree        | local main commit  |
      | existing | local, worktree | existing commit    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the initial branches and lineage exist now
    And the initial commits exist now
