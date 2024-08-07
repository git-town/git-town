Feature: don't sync tags while prepending

  Background:
    Given a Git repo with origin
    And the branch
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | old  | feature | main   | local, origin |
    And the tags
      | NAME       | LOCATION |
      | local-tag  | local    |
      | origin-tag | origin   |
    And the current branch is "old"
    And Git Town setting "sync-tags" is "false"
    When I run "git-town prepend new"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                             |
      | old    | git fetch --prune --no-tags         |
      |        | git checkout main                   |
      | main   | git rebase origin/main              |
      |        | git checkout old                    |
      | old    | git merge --no-edit --ff origin/old |
      |        | git merge --no-edit --ff main       |
      |        | git checkout -b new main            |
    And the initial tags exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND           |
      | new    | git checkout old  |
      | old    | git branch -D new |
    And the initial commits exist
    And the initial lineage exists
    And the initial tags exist now
