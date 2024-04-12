Feature: compress the commits on an entire stack when at the stack root

  Background:
    Given feature branch "alpha" with these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | alpha 1 | alpha_1   | alpha 1      |
      |               | alpha 2 | alpha_2   | alpha 2      |
      |               | alpha 3 | alpha_3   | alpha 3      |
    And feature branch "beta" as a child of "alpha" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | beta 1  | beta_1    | beta 1       |
      |               | beta 2  | beta_2    | beta 2       |
      |               | beta 3  | beta_3    | beta 3       |
    And feature branch "gamma" as a child of "beta" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | gamma 1 | gamma_1   | gamma 1      |
      |               | gamma 2 | gamma_2   | gamma 2      |
      |               | gamma 3 | gamma_3   | gamma 3      |
    And the current branch is "gamma"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | gamma  | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash                                       |
      |        | git checkout alpha                              |
      | alpha  | git reset --soft main                           |
      |        | git commit -m "alpha 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --soft alpha                          |
      |        | git commit -m "beta 1"                          |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma                              |
      | gamma  | git reset --soft beta                           |
      |        | git commit -m "gamma 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "gamma"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE |
      | alpha  | local, origin | alpha 1 |
      | beta   | local, origin | alpha 1 |
      |        |               | beta 1  |
      | gamma  | local, origin | alpha 1 |
      |        |               | beta 1  |
      |        |               | gamma 1 |
    And file "alpha_1" still has content "alpha 1"
    And file "alpha_2" still has content "alpha 2"
    And file "alpha_3" still has content "alpha 3"
    And file "beta_1" still has content "beta 1"
    And file "beta_2" still has content "beta 2"
    And file "beta_3" still has content "beta 3"
    And file "gamma_1" still has content "gamma 1"
    And file "gamma_2" still has content "gamma 2"
    And file "gamma_3" still has content "gamma 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | gamma  | git add -A                                      |
      |        | git stash                                       |
      |        | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'alpha 3' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout beta                               |
      | beta   | git reset --hard {{ sha 'beta 3' }}             |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout gamma                              |
      | gamma  | git reset --hard {{ sha 'gamma 3' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And the current branch is still "gamma"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
