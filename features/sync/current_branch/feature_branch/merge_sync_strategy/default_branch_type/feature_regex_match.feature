@smoke
Feature: a unknown branch type is set, the feature-regex matches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE   | PARENT | LOCATIONS     |
      | my-branch | (none) | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION | MESSAGE                 |
      | main      | local    | local main commit       |
      |           | origin   | origin main commit      |
      | my-branch | local    | local my-branch commit  |
      |           | origin   | origin my-branch commit |
    And local Git setting "git-town.feature-regex" is "my-.*"
    And local Git setting "git-town.unknown-branch-type" is "observed"
    And the current branch is "my-branch"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | my-branch | git fetch --prune --tags                          |
      |           | git checkout main                                 |
      | main      | git -c rebase.updateRefs=false rebase origin/main |
      |           | git push                                          |
      |           | git checkout my-branch                            |
      | my-branch | git merge --no-edit --ff main                     |
      |           | git merge --no-edit --ff origin/my-branch         |
      |           | git push                                          |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                                                        |
      | main      | local, origin | origin main commit                                             |
      |           |               | local main commit                                              |
      | my-branch | local, origin | local my-branch commit                                         |
      |           |               | Merge branch 'main' into my-branch                             |
      |           |               | origin my-branch commit                                        |
      |           |               | Merge remote-tracking branch 'origin/my-branch' into my-branch |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                                    |
      | my-branch | git reset --hard {{ sha 'local my-branch commit' }}                                        |
      |           | git push --force-with-lease origin {{ sha-in-origin 'origin my-branch commit' }}:my-branch |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                 |
      | main      | local, origin | origin main commit      |
      |           |               | local main commit       |
      | my-branch | local         | local my-branch commit  |
      |           | origin        | origin my-branch commit |
