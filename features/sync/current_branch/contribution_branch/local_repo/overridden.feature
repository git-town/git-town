Feature: sync the current branch which has an override to "contribution"

  Background:
    Given a local Git repo
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS |
      | contribution | contribution | main   | local     |
    And the commits
      | BRANCH       | LOCATION | MESSAGE             |
      | main         | local    | main commit         |
      | contribution | local    | contribution commit |
    And the current branch is "contribution"
    And Git setting "git-town-branch.contribution.branch-type" is "feature"
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town prints:
      """
      xxx
      """
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git push                                |
    And all branches are now synchronized
    And the current branch is still "contribution"
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
    Then Git Town runs the commands
      | BRANCH       | COMMAND       |
      | contribution | git add -A    |
      |              | git stash     |
      |              | git stash pop |
    And the current branch is still "contribution"
    And the initial commits exist now
    And the initial branches and lineage exist now
