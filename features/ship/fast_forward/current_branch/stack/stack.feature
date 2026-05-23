Feature: ship a parent branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME     | FILE CONTENT      |
      | feature-1 | local, origin | feature-1 commit | feature-1.txt | feature-1 content |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME     | FILE CONTENT      |
      | feature-2 | local, origin | feature-2 commit | feature-2.txt | feature-2 content |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-3 | feature | feature-2 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME     | FILE CONTENT      |
      | feature-3 | local, origin | feature-3 commit | feature-3.txt | feature-3 content |
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And the current branch is "feature-1"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                       |
      | feature-1 | git fetch --prune --tags      |
      |           | git checkout main             |
      | main      | git merge --ff-only feature-1 |
      |           | git push                      |
      |           | git push origin :feature-1    |
      |           | git branch -D feature-1       |
    And Git Town prints:
      """
      branch feature-2 is now a child of main
      """
    And this lineage exists now
      """
      main
        feature-2
          feature-3
      """
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, origin | feature-1 commit |
      | feature-2 | local, origin | feature-2 commit |
      | feature-3 | local, origin | feature-3 commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git branch feature-1 {{ sha 'feature-1 commit' }} |
      |        | git push -u origin feature-1                      |
      |        | git checkout feature-1                            |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, origin | feature-1 commit |
      | feature-2 | local, origin | feature-2 commit |
      | feature-3 | local, origin | feature-3 commit |

  Scenario: ship the second branch
    Given the current branch is "feature-2"
    When I run "git-town ship"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                       |
      | feature-2 | git fetch --prune --tags      |
      |           | git checkout main             |
      | main      | git merge --ff-only feature-2 |
      |           | git push                      |
      |           | git push origin :feature-2    |
      |           | git branch -D feature-2       |
    And Git Town prints:
      """
      branch feature-3 is now a child of main
      """
    And this lineage exists now
      """
      main
        feature-3
      """
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, origin | feature-1 commit |
      |           |               | feature-2 commit |
      | feature-3 | local, origin | feature-3 commit |

  Scenario: ship all remaining branches
    Given the current branch is "feature-2"
    And I run "git-town ship"
    And the current branch is "feature-3"
    When I run "git-town ship"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                       |
      | feature-3 | git fetch --prune --tags      |
      |           | git checkout main             |
      | main      | git merge --ff-only feature-3 |
      |           | git push                      |
      |           | git push origin :feature-3    |
      |           | git branch -D feature-3       |
    And Git Town prints:
      """
      deleted branch feature-3
      """
    And no lineage exists now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE          |
      | main   | local, origin | feature-1 commit |
      |        |               | feature-2 commit |
      |        |               | feature-3 commit |
