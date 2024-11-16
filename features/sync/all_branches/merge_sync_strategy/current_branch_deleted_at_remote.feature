@this
Feature: sync a branch whose tracking branch was shipped

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE      | PARENT | LOCATIONS     |
      | hooks | prototype | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | hooks  | local, origin | hooks commit |
    And origin ships the "hooks" branch
    And the current branch is "hooks"
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
    And the initial branches and lineage exist now
