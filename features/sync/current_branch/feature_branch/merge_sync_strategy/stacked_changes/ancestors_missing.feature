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

  @this
  Scenario:
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | gamma  | git fetch --prune --tags              |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git merge --no-edit --ff origin/main  |
      |        | git push                              |
      |        | git checkout gamma                    |
      | gamma  | git merge --no-edit --ff origin/gamma |
      |        | git merge --no-edit --ff origin/beta  |
      |        | git merge --no-edit --ff alpha        |
      |        | git push                              |
    And all branches are now synchronized
    And the current branch is still "gamma"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                |
      | main   | origin        | origin main commit                                     |
      | alpha  | local, origin | local alpha commit                                     |
      |        |               | origin alpha commit                                    |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha |
      |        |               | origin main commit                                     |
      |        |               | Merge remote-tracking branch 'origin/main' into alpha  |
      | beta   | origin        | origin beta commit                                     |
      | gamma  | local, origin | local gamma commit                                     |
      |        |               | origin gamma commit                                    |
      |        |               | Merge remote-tracking branch 'origin/gamma' into gamma |
      |        |               | origin beta commit                                     |
      |        |               | Merge remote-tracking branch 'origin/beta' into gamma  |
      |        |               | local alpha commit                                     |
      |        |               | origin alpha commit                                    |
      |        |               | Merge remote-tracking branch 'origin/alpha' into alpha |
      |        |               | origin main commit                                     |
      |        |               | Merge remote-tracking branch 'origin/main' into alpha  |
      |        |               | Merge branch 'alpha' into gamma                        |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                                         |
      | child  | git reset --hard {{ sha-before-run 'local child commit' }}                                      |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin child commit' }}:child   |
      |        | git checkout parent                                                                             |
      | parent | git reset --hard {{ sha-before-run 'local parent commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin parent commit' }}:parent |
      |        | git checkout child                                                                              |
