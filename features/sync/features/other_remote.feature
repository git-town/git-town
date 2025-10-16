Feature: ignores other Git remotes

  Background:
    Given a Git repo with origin
    And a remote "other" pointing to "git@foo.com:bar/baz.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | TYPE     | COMMAND                                 |
      | feature | frontend | git fetch --prune --tags                |
      |         | frontend | git merge --no-edit --ff origin/feature |
      |         | frontend | git push                                |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                |
      | feature | git reset --hard {{ sha 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'origin feature commit' }}:feature |
    And the initial branches and lineage exist now
    And the initial commits exist now
