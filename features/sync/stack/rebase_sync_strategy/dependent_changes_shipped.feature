Feature: shipped the head branch of a synced stack with dependent changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE            | FILE NAME | FILE CONTENT  |
      | alpha  | local    | local alpha commit | file      | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE           | FILE NAME | FILE CONTENT |
      | beta   | local    | local beta commit | file      | beta content |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME | FILE CONTENT  |
      | alpha  | local, origin | local alpha commit | file      | alpha content |
      | beta   | local, origin | local alpha commit | file      | alpha content |
      |        |               | local beta commit  | file      | beta content  |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the current branch is "beta"
    And I ran "git-town sync"
    And it runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main --no-update-refs         |
      |        | git checkout alpha                              |
      | alpha  | git rebase main --no-update-refs                |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git rebase alpha --no-update-refs               |
      |        | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME | FILE CONTENT  |
      | alpha  | local, origin | local alpha commit | file      | alpha content |
      | beta   | local, origin | local alpha commit | file      | alpha content |
      |        |               | local beta commit  | file      | beta content  |
    And origin ships the "alpha" branch
    # And inspect the repo
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git fetch --prune --tags                        |
      |        | git checkout main                               |
      | main   | git rebase origin/main --no-update-refs         |
      |        | git branch -D alpha                             |
      |        | git checkout beta                               |
      | beta   | git rebase main --no-update-refs                |
      |        | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | local, origin | origin main commit  |
      |        |               | local main commit   |
      | beta   | local, origin | origin beta commit  |
      |        |               | origin alpha commit |
      |        |               | origin main commit  |
      |        |               | local main commit   |
      |        |               | local alpha commit  |
      |        |               | local beta commit   |
      | alpha  | local, origin | origin alpha commit |
      |        |               | origin main commit  |
      |        |               | local main commit   |
      |        |               | local alpha commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                                       |
      | beta   | git reset --hard {{ sha-before-run 'local beta commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin beta commit' }}:beta   |
      |        | git checkout alpha                                                                            |
      | alpha  | git reset --hard {{ sha-before-run 'local alpha commit' }}                                    |
      |        | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin alpha commit' }}:alpha |
      |        | git checkout beta                                                                             |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE             |
      | main   | local, origin | origin main commit  |
      |        |               | local main commit   |
      | beta   | local         | local beta commit   |
      |        | origin        | origin beta commit  |
      | alpha  | local         | local alpha commit  |
      |        | origin        | origin alpha commit |
    And the initial branches and lineage exist now
