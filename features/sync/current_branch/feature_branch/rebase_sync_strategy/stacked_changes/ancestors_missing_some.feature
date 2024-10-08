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
    And Git Town setting "sync-feature-strategy" is "rebase"
    And I ran "git branch -d main"
    And I ran "git branch -d beta"
    When I run "git-town sync"

  Scenario:
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | gamma  | git fetch --prune --tags                        |
      |        | git checkout alpha                              |
      | alpha  | git rebase origin/main                          |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/alpha                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma                              |
      | gamma  | git rebase origin/beta                          |
      |        | git rebase alpha --no-update-refs               |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/gamma                         |
      |        | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "gamma"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | origin        | origin main commit  |
      | alpha  | local, origin | origin alpha commit |
      |        |               | origin main commit  |
      |        |               | local alpha commit  |
      | beta   | origin        | origin beta commit  |
      | gamma  | local, origin | origin gamma commit |
      |        |               | origin alpha commit |
      |        |               | origin main commit  |
      |        |               | local alpha commit  |
      |        |               | origin beta commit  |
      |        |               | local gamma commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
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
