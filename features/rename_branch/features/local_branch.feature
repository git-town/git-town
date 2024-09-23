Feature: rename a local-only branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS |
      | old  | feature | main   | local     |
    And the current branch is "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | old    | local         | old commit  |
    When I run "git-town rename-branch new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                   |
      | old    | git fetch --prune --tags  |
      |        | git branch --move old new |
      |        | git checkout new          |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE     |
      | main   | local, origin | main commit |
      | new    | local         | old commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | new    | git branch old {{ sha 'old commit' }} |
      |        | git checkout old                      |
      | old    | git branch -D new                     |
    And the current branch is now "old"
    And the initial branches and lineage exist now
