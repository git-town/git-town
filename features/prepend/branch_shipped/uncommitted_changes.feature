Feature: prepend a branch to a branch that was shipped at the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin ships the "child" branch using the "squash-merge" ship-strategy
    And the current branch is "child"
    And an uncommitted file
    When I run "git-town prepend new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | child  | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout -b new parent  |
      | new    | git stash pop               |
      |        | git restore --staged .      |
    And Git Town prints:
      """
      branch "new" is now a child of "parent"
      """
    And Git Town prints:
      """
      branch "child" is now a child of "new"
      """
    And this lineage exists now
      """
      main
        parent
          new
            child
      """
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | main, child, new, parent |
      | origin     | main, parent             |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | new    | git add -A                  |
      |        | git stash -m "Git Town WIP" |
      |        | git checkout child          |
      | child  | git branch -D new           |
      |        | git stash pop               |
      |        | git restore --staged .      |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY | BRANCHES            |
      | local      | main, child, parent |
      | origin     | main, parent        |
    And the uncommitted file still exists
