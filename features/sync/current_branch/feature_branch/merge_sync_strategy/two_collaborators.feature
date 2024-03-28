Feature: collaborative feature branch syncing

  Scenario: two computers work on a feature branch
    Given a coworker clones the repository
    And the current branch is a feature branch "feature"
    And the coworker fetches updates
    And the coworker sets the parent branch of "feature" as "main"
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | feature | local    | my commit       |
      |         | coworker | coworker commit |
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         |
      | feature | local, origin | my commit       |
      |         | coworker      | coworker commit |
    And all branches are now synchronized

    Given the coworker is on the "feature" branch
    When the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE                                                    |
      | feature | local, coworker, origin | my commit                                                  |
      |         | coworker, origin        | coworker commit                                            |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature |

    Given the current branch is "feature"
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE                                                    |
      | feature | local, coworker, origin | coworker commit                                            |
      |         |                         | my commit                                                  |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature |
