Feature: sync a feature branch with new commits on the tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT |
      | alpha  | local, origin | alpha commit | alpha_file | content 1    |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | beta_file | content 2    |
      | beta   | origin        | new commit  | beta_file | content 3    |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And the current branch is "beta"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                              |
      | beta   | git fetch --prune --tags             |
      |        | git checkout alpha                   |
      | alpha  | git checkout beta                    |
      | beta   | git merge --no-edit --ff origin/beta |
      |        | git reset --soft alpha --            |
      |        | git commit -m "beta commit"          |
      |        | git push --force-with-lease          |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT |
      | alpha  | local, origin | alpha commit | alpha_file | content 1    |
      | beta   | local, origin | beta commit  | beta_file  | content 3    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                  |
      | beta   | git reset --hard {{ sha-initial 'beta commit' }}                         |
      |        | git push --force-with-lease origin {{ sha-in-origin 'new commit' }}:beta |
    And the initial branches and lineage exist now
    And the initial commits exist now
