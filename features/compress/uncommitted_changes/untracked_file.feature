Feature: compress the commits on an entire stack when at the stack root

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | alpha  | local, origin | alpha 1 | alpha_1   | alpha 1      |
      |        |               | alpha 2 | alpha_2   | alpha 2      |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta 1  | beta_1    | beta 1       |
      |        |               | beta 2  | beta_2    | beta 2       |
    And the current branch is "alpha"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash -m "Git Town WIP"                     |
      |        | git reset --soft main --                        |
      |        | git commit -m "alpha 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --soft alpha --                       |
      |        | git commit -m "beta 1"                          |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout alpha                              |
      | alpha  | git stash pop                                   |
      |        | git restore --staged .                          |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE |
      | alpha  | local, origin | alpha 1 |
      | beta   | local, origin | beta 1  |
    And file "alpha_1" still has content "alpha 1"
    And file "alpha_2" still has content "alpha 2"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | alpha  | git add -A                                      |
      |        | git stash -m "Git Town WIP"                     |
      |        | git reset --hard {{ sha 'alpha 2' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'beta 2' }}             |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout alpha                              |
      | alpha  | git stash pop                                   |
      |        | git restore --staged .                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
