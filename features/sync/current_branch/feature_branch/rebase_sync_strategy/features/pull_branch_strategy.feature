Feature: with pull-branch-strategy set to "merge"

  Background:
    Given setting "sync-strategy" is "rebase"
    And setting "pull-branch-strategy" is "merge"
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --no-edit origin/main |
      |         | git push                        |
      |         | git checkout feature            |
      | feature | git rebase origin/feature       |
      |         | git rebase main                 |
      |         | git push --force-with-lease     |
    And all branches are now synchronized
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                                    |
      | main    | local, origin | local main commit                          |
      |         |               | origin main commit                         |
      |         |               | Merge remote-tracking branch 'origin/main' |
      | feature | local, origin | local main commit                          |
      |         |               | origin main commit                         |
      |         |               | Merge remote-tracking branch 'origin/main' |
      |         |               | origin feature commit                      |
      |         |               | local feature commit                       |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                                           |
      | feature | git reset --hard {{ sha-before-run 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin feature commit' }}:feature |
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                                    |
      | main    | local, origin | local main commit                          |
      |         |               | origin main commit                         |
      |         |               | Merge remote-tracking branch 'origin/main' |
      | feature | local         | local feature commit                       |
      |         | origin        | origin feature commit                      |
    And the initial branches and hierarchy exist
