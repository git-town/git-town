@smoke
Feature: a default branch type is set, the feature-regex matches

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
    And the current branch is "my-branch"
    And local Git Town setting "feature-regex" is "my-.*"
    And local Git Town setting "default-branch-type" is "observed"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                                   |
      | my-branch | git fetch --prune --tags                  |
      |           | git checkout main                         |
      | main      | git rebase origin/main --no-update-refs   |
      |           | git push                                  |
      |           | git checkout my-branch                    |
      | my-branch | git merge --no-edit --ff main             |
      |           | git merge --no-edit --ff origin/my-branch |
      |           | git push                                  |
    And all branches are now synchronized
    And the current branch is still "my-branch"
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
    Then it runs the commands
      | BRANCH    | COMMAND                                                                                    |
      | my-branch | git reset --hard {{ sha 'local my-branch commit' }}                                        |
      |           | git push --force-with-lease origin {{ sha-in-origin 'origin my-branch commit' }}:my-branch |
    And the current branch is still "my-branch"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE                 |
      | main      | local, origin | origin main commit      |
      |           |               | local main commit       |
      | my-branch | local         | local my-branch commit  |
      |           | origin        | origin my-branch commit |
    And the initial branches and lineage exist now
