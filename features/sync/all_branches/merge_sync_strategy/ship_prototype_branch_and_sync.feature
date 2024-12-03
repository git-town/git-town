Feature: end-to-end workflow of creating a prototype branch, shipping, and pruning it

  This test reproduces the bug in https://github.com/git-town/git-town/issues/4222.

  Background:
    Given a Git repo with origin
    And Git Town setting "default-branch-type" is "prototype"
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    And I ran "git-town hack hooks"
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | hooks  | local    | hooks commit |
    And the current branch is "hooks"
    And I ran "git-town propose"
    And Git Town prints:
      """
      branch 'hooks' set up to track 'origin/hooks'
      """
    And Git Town prints:
      """
      branch "hooks" is no longer a prototype branch
      """
    And origin ships the "hooks" branch using the "squash-merge" ship-strategy
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | hooks  | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git branch -D hooks                     |
      |        | git push --tags                         |
    And Git Town prints:
      """
      deleted branch "hooks"
      """
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git reset --hard {{ sha 'initial commit' }} |
      |        | git branch hooks {{ sha 'hooks commit' }}   |
      |        | git checkout hooks                          |
    And the current branch is now "hooks"
    And the branches are now
      | REPOSITORY | BRANCHES    |
      | local      | main, hooks |
      | origin     | main        |
    And this lineage exists now
      | BRANCH | PARENT |
      | hooks  | main   |
