@messyoutput
Feature: sync the current branch which has a branch-type override

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | LOCATIONS     |
      | contribution | contribution | local, origin |
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | main         | local, origin | main commit         |
      | contribution | local         | contribution commit |
    And the current branch is "contribution"
    And I ran "git-town feature"
    When I run "git-town sync" and enter into the dialog:
      | DIALOG                           | KEYS  |
      | parent branch for "contribution" | enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                      |
      | contribution | git fetch --prune --tags                     |
      |              | git merge --no-edit --ff main                |
      |              | git merge --no-edit --ff origin/contribution |
      |              | git push                                     |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH       | LOCATION      | MESSAGE                               |
      | main         | local, origin | main commit                           |
      | contribution | local, origin | contribution commit                   |
      |              |               | Merge branch 'main' into contribution |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                                                    |
      | contribution | git reset --hard {{ sha 'contribution commit' }}                           |
      |              | git push --force-with-lease origin {{ sha 'initial commit' }}:contribution |
    And the initial branches and lineage exist now
    And the initial commits exist now
