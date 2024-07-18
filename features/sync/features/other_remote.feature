Feature: ignores other Git remotes

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    Given the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And a remote "other" pointing to "git@foo.com:bar/baz.git"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                 |
      | feature | frontend | git fetch --prune --tags                |
      |         | frontend | git checkout main                       |
      | main    | frontend | git rebase origin/main                  |
      |         | frontend | git checkout feature                    |
      | feature | frontend | git merge --no-edit --ff origin/feature |
      |         | frontend | git merge --no-edit --ff main           |
      |         | frontend | git push                                |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
