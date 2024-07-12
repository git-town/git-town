Feature: syncing a top-level feature branch using --no-push

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync --no-push"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, origin | origin main commit                                         |
      |         | local         | local main commit                                          |
      | feature | local         | local feature commit                                       |
      |         | local, origin | origin feature commit                                      |
      |         | local         | Merge remote-tracking branch 'origin/feature' into feature |
      |         |               | origin main commit                                         |
      |         |               | local main commit                                          |
      |         |               | Merge branch 'main' into feature                           |
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
    And it runs no commands
