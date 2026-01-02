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
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "gamma"
    And I ran "git branch -d beta"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                       |
      | gamma  | git fetch --prune --tags                                                      |
      |        | git checkout alpha                                                            |
      | alpha  | git push --force-with-lease --force-if-includes                               |
      |        | git -c rebase.updateRefs=false rebase origin/alpha                            |
      |        | git -c rebase.updateRefs=false rebase origin/main                             |
      |        | git push --force-with-lease --force-if-includes                               |
      |        | git checkout gamma                                                            |
      | gamma  | git push --force-with-lease --force-if-includes                               |
      |        | git -c rebase.updateRefs=false rebase origin/gamma                            |
      |        | git -c rebase.updateRefs=false rebase origin/beta                             |
      |        | git -c rebase.updateRefs=false rebase --onto alpha {{ sha 'initial commit' }} |
      |        | git push --force-with-lease --force-if-includes                               |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | origin        | origin main commit  |
      | alpha  | local, origin | origin alpha commit |
      |        |               | local alpha commit  |
      | beta   | origin        | origin beta commit  |
      | gamma  | local, origin | origin main commit  |
      |        |               | origin alpha commit |
      |        |               | local alpha commit  |
      |        |               | origin beta commit  |
      |        |               | origin gamma commit |
      |        |               | local gamma commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                    |
      | gamma  | git checkout alpha                                                                         |
      | alpha  | git reset --hard {{ sha-initial 'local alpha commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin alpha commit' }}:alpha |
      |        | git checkout gamma                                                                         |
      | gamma  | git reset --hard {{ sha-initial 'local gamma commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin gamma commit' }}:gamma |
    And the initial lineage exists now
    And the branches are now
      | REPOSITORY | BRANCHES                 |
      | local      | alpha, gamma             |
      | origin     | main, alpha, beta, gamma |
    And the initial commits exist now
