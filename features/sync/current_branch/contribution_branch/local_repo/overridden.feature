Feature: sync the current branch which has a branch-type override

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | prototype | prototype | main   | local     |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, origin | main commit      |
      | prototype | local         | prototype commit |
    And the current branch is "prototype"
    And Git setting "git-town-branch.prototype.branchtype" is "feature"
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                 |
      | prototype | git fetch --prune --tags                |
      |           | git checkout main                       |
      | main      | git rebase origin/main --no-update-refs |
      |           | git checkout prototype                  |
      | prototype | git merge --no-edit --ff main           |
      |           | git push -u origin prototype            |
    And all branches are now synchronized
    And the current branch is still "prototype"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                            |
      | main      | local, origin | main commit                        |
      | prototype | local, origin | prototype commit                   |
      |           |               | Merge branch 'main' into prototype |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND       |
      | prototype | git add -A    |
      |           | git stash     |
      |           | git stash pop |
    And the current branch is still "prototype"
    And the initial commits exist now
    And the initial branches and lineage exist now
