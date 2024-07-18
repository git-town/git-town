Feature: hack a new branch while the main branch is active in another worktree

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the current branch is "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE            |
      | main     | origin   | origin main commit |
      |          | local    | local main commit  |
      | existing | local    | existing commit    |
    And branch "main" is active in another worktree
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git fetch --prune --tags |
      |          | git checkout -b new main |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION        | MESSAGE            |
      | main     | origin          | origin main commit |
      |          | worktree        | local main commit  |
      | existing | local, worktree | existing commit    |
      | new      | local           | local main commit  |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the current branch is now "existing"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE            |
      | main     | origin   | origin main commit |
      |          | worktree | local main commit  |
      | existing | local    | existing commit    |
    And the initial branches and lineage exist
