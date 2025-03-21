Feature: prune enabled via CLI

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town sync --prune"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git push --no-verify                    |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push --no-verify                    |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                                                    |
      | main    | local, origin | origin main commit                                         |
      |         |               | local main commit                                          |
      | feature | local, origin | local feature commit                                       |
      |         |               | Merge branch 'main' into feature                           |
      |         |               | origin feature commit                                      |
      |         |               | Merge remote-tracking branch 'origin/feature' into feature |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
