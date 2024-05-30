Feature: other Git remotes exist

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And a remote "other" pointing to "git@foo.com:bar/baz.git"
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                 |
      | feature | frontend | git fetch --prune --tags                |
      |         | frontend | git checkout main                       |
      | main    | frontend | git rebase origin/main                  |
      |         | frontend | git push                                |
      |         | frontend | git checkout feature                    |
      | feature | frontend | git merge --no-edit --ff origin/feature |
      |         | frontend | git merge --no-edit --ff main           |
      |         | frontend | git push                                |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                               |
      | feature | git push --force-with-lease origin {{ sha 'initial commit' }}:feature |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
