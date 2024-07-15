Feature: syncing a top-level feature branch using --no-push

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town sync --no-push"

  @debug @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                   |
      | feature | git fetch --prune --tags  |
      |         | git checkout main         |
      | main    | git rebase origin/main    |
      |         | git checkout feature      |
      | feature | git rebase main           |
      |         | git rebase origin/feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | origin main commit    |
      |         | local         | local main commit     |
      | feature | local, origin | origin feature commit |
      |         |               | origin main commit    |
      |         |               | local main commit     |
      |         |               | local feature commit  |
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git reset --hard {{ sha 'local feature commit' }}                            |
      |         | git push --force-with-lease origin {{ sha 'origin feature commit' }}:feature |
      |         | git checkout main                                                            |
      | main    | git reset --hard {{ sha 'local main commit' }}                               |
      |         | git checkout feature                                                         |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
