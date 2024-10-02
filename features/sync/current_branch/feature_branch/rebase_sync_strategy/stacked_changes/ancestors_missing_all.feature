@this
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
    And Git Town setting "sync-feature-strategy" is "rebase"
    And I ran "git branch -d main"
    And I ran "git branch -d alpha"
    And I ran "git branch -d beta"
    When I run "git-town sync"

  Scenario:
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | gamma  | git fetch --prune --tags                        |
      |        | git rebase origin/beta                          |
      |        | git rebase origin/alpha                         |
      |        | git rebase origin/main                          |
      |        | git push --force-with-lease --force-if-includes |
      |        | git rebase origin/gamma                         |
      |        | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "gamma"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | origin        | origin main commit  |
      | alpha  | origin        | origin alpha commit |
      | beta   | origin        | origin beta commit  |
      | gamma  | local, origin | origin gamma commit |
      |        |               | origin main commit  |
      |        |               | origin alpha commit |
      |        |               | origin beta commit  |
      |        |               | local gamma commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                                       |
      | gamma  | git reset --hard {{ sha-before-run 'local gamma commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin gamma commit' }}:gamma |
    And the initial lineage exists now
    And the initial commits exist now
    And these branches exist now
      | REPOSITORY | BRANCHES                 |
      | local      | gamma                    |
      | origin     | main, alpha, beta, gamma |
