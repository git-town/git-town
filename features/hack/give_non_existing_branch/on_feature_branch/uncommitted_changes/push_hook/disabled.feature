Feature: auto-push the new branch without running Git push hooks

  Background:
    Given a Git repo with origin
    And Git Town setting "push-new-branches" is "true"
    And Git Town setting "push-hook" is "false"
    And the commits
      | BRANCH | LOCATION | MESSAGE       |
      | main   | origin   | origin commit |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                            |
      | main   | git add -A                         |
      |        | git stash                          |
      |        | git checkout -b new                |
      | new    | git push --no-verify -u origin new |
      |        | git stash pop                      |
    And the current branch is now "new"
    And the initial commits exist now
    And this lineage exists now
      | BRANCH | PARENT |
      | new    | main   |
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | new    | git add -A           |
      |        | git stash            |
      |        | git checkout main    |
      | main   | git branch -D new    |
      |        | git push origin :new |
      |        | git stash pop        |
    And the current branch is now "main"
    And the uncommitted file still exists
    And the initial commits exist now
    And no lineage exists now
