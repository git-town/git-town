Feature: stacked changes where all ancestor branches aren't local

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
      | alpha  | origin   | origin alpha commit |
      | beta   | origin   | origin beta commit  |
      | gamma  | local    | local gamma commit  |
      |        | origin   | origin gamma commit |
    And the current branch is "gamma"
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And I ran "git branch -d main"
    And I ran "git branch -d alpha"
    And I ran "git branch -d beta"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | gamma  | git fetch --prune --tags              |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git merge --no-edit --ff origin/alpha |
      |        | git merge --no-edit --ff origin/main  |
      |        | git merge --no-edit --ff origin/gamma |
      |        | git reset --soft origin/beta          |
      |        | git commit -m "local gamma commit"    |
      |        | git push --force-with-lease           |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | origin        | origin main commit  |
      | alpha  | origin        | origin alpha commit |
      | beta   | origin        | origin beta commit  |
      | gamma  | local, origin | local gamma commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                    |
      | gamma  | git reset --hard {{ sha-initial 'local gamma commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin gamma commit' }}:gamma |
    And the initial lineage exists now
    And the initial commits exist now
    And these branches exist now
      | REPOSITORY | BRANCHES                 |
      | local      | gamma                    |
      | origin     | main, alpha, beta, gamma |
