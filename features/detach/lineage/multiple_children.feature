Feature: detach a branch branch with multiple children

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | alpha  | feature | main   | local, origin |
      | beta   | feature | alpha  | local, origin |
      | gamma1 | feature | beta   | local, origin |
      | gamma2 | feature | beta   | local, origin |
      | delta  | feature | gamma2 | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME   | FILE CONTENT   |
      | main   | local, origin | main commit   | main-file   | main content   |
      | alpha  | local, origin | alpha commit  | alpha-file  | alpha content  |
      | beta   | local, origin | beta commit   | beta-file   | beta content   |
      | gamma1 | local, origin | gamma1 commit | gamma1-file | gamma1 content |
      | gamma2 | local, origin | gamma2 commit | gamma2-file | gamma2 content |
      | delta  | local, origin | delta commit  | delta-file  | delta content  |
    And the current branch is "beta"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                  |
      | beta   | git fetch --prune --tags                                 |
      |        | Finding proposal from beta into alpha ... none           |
      |        | Finding proposal from gamma1 into beta ... none          |
      |        | Finding proposal from gamma2 into beta ... none          |
      |        | git checkout gamma1                                      |
      | gamma1 | git pull                                                 |
      |        | git -c rebase.updateRefs=false rebase --onto alpha beta  |
      |        | git push --force-with-lease                              |
      |        | git checkout gamma2                                      |
      | gamma2 | git pull                                                 |
      |        | git -c rebase.updateRefs=false rebase --onto alpha beta  |
      |        | git push --force-with-lease                              |
      |        | git checkout delta                                       |
      | delta  | git pull                                                 |
      |        | git -c rebase.updateRefs=false rebase --onto gamma2 beta |
      |        | git push --force-with-lease                              |
      |        | git checkout beta                                        |
      | beta   | git -c rebase.updateRefs=false rebase --onto main alpha  |
      |        | git push --force-with-lease --force-if-includes          |
    And this lineage exists now
      """
      main
        alpha
          gamma1
          gamma2
            delta
        beta
      """
    And the branches are now
      | REPOSITORY    | BRANCHES                                 |
      | local, origin | main, alpha, beta, delta, gamma1, gamma2 |
    And no uncommitted files exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git reset --hard {{ sha 'beta commit' }}        |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout delta                              |
      | delta  | git reset --hard {{ sha 'delta commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma1                             |
      | gamma1 | git reset --hard {{ sha 'gamma1 commit' }}      |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma2                             |
      | gamma2 | git reset --hard {{ sha 'gamma2 commit' }}      |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
    And the initial branches and lineage exist now
    And the initial commits exist now
