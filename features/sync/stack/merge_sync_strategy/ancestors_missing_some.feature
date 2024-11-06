Feature: stacked changes where an ancestor branch isn't local

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
      | gamma | feature | beta   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE             |
      | main   | origin   | origin main commit  |
      | alpha  | local    | local alpha commit  |
      | alpha  | origin   | origin alpha commit |
      | beta   | origin   | origin beta commit  |
      | gamma  | local    | local gamma commit  |
      |        | origin   | origin gamma commit |
    And the current branch is "gamma"
    And I ran "git branch -d main"
    And I ran "git branch -d beta"
    When I run "git-town sync"

  Scenario:
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | gamma  | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/main  |
      |        | git merge --no-edit --ff origin/alpha |
      |        | git push                              |
      |        | git checkout gamma                    |
      | gamma  | git merge --no-edit --ff origin/beta  |
      |        | git merge --no-edit --ff alpha        |
      |        | git merge --no-edit --ff origin/gamma |
      |        | git push                              |
    And all branches are now synchronized
    And the current branch is still "gamma"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                |
      | main   | origin        | origin main commit                                     |
      | alpha  | local, origin | local alpha commit                                     |
      |        |               | origin main commit                                     |
      |        |               | Merge remote-tracking branch 'origin/main' into alpha  |
      |        |               | origin alpha commit                                    |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha |
      | beta   | origin        | origin beta commit                                     |
      | gamma  | local, origin | local gamma commit                                     |
      |        |               | origin beta commit                                     |
      |        |               | Merge remote-tracking branch 'origin/beta' into gamma  |
      |        |               | local alpha commit                                     |
      |        |               | origin main commit                                     |
      |        |               | Merge remote-tracking branch 'origin/main' into alpha  |
      |        |               | origin alpha commit                                    |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha |
      |        |               | Merge branch 'alpha' into gamma                        |
      |        |               | origin gamma commit                                    |
      |        |               | Merge remote-tracking branch 'origin/gamma' into gamma |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                       |
      | gamma  | git checkout alpha                                                                            |
      | alpha  | git reset --hard {{ sha-before-run 'local alpha commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin alpha commit' }}:alpha |
      |        | git checkout gamma                                                                            |
      | gamma  | git reset --hard {{ sha-before-run 'local gamma commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin gamma commit' }}:gamma |
    And the initial lineage exists now
    And the initial commits exist now
    And these branches exist now
      | REPOSITORY | BRANCHES                 |
      | local      | alpha, gamma             |
      | origin     | main, alpha, beta, gamma |
