Feature: syncing a local feature branch using --no-push

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      |         | origin   | origin main commit   |
      | feature | local    | local feature commit |
    When I run "git-town sync --no-push"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git fetch --prune --tags      |
      |         | git checkout main             |
      | main    | git rebase origin/main        |
      |         | git checkout feature          |
      | feature | git merge --no-edit --ff main |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | origin main commit               |
      |         | local         | local main commit                |
      | feature | local         | local feature commit             |
      |         |               | origin main commit               |
      |         |               | local main commit                |
      |         |               | Merge branch 'main' into feature |
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git reset --hard {{ sha 'local feature commit' }} |
      |         | git checkout main                                 |
      | main    | git reset --hard {{ sha 'local main commit' }}    |
      |         | git checkout feature                              |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
