Feature: sync the current branch which has a branch-type override

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS |
      | contribution | contribution | local     |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | main         | local, origin | main commit         |
      | contribution | local         | contribution commit |
    And the current branch is "contribution"
    And I ran "git-town hack"
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                          | KEYS  |
      | parent branch of "contribution" | enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                 |
      | contribution | git fetch --prune --tags                |
      |              | git checkout main                       |
      | main         | git rebase origin/main --no-update-refs |
      |              | git checkout contribution               |
      | contribution | git merge --no-edit --ff main           |
      |              | git push -u origin contribution         |
    And all branches are now synchronized
    And the current branch is still "contribution"
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE                               |
      | main         | local, origin | main commit                           |
      | contribution | local, origin | contribution commit                   |
      |              |               | Merge branch 'main' into contribution |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                          |
      | contribution | git reset --hard {{ sha 'contribution commit' }} |
      |              | git push origin :contribution                    |
    And the current branch is still "contribution"
    And the initial commits exist now
    And the initial branches and lineage exist now
