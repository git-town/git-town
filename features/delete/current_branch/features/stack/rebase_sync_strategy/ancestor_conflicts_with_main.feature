Feature: deleting a branch that conflicts with the main branch

  Background:
    Given a Git repo with origin
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | main   | local, origin | main commit | file      | main content |
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME | FILE CONTENT |
      | feature-1 | local, origin | feature-1 commit | file      | content 1    |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME | FILE CONTENT |
      | feature-2 | local, origin | feature-2 commit | file      | content 2    |
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-3 | feature | feature-2 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME | FILE CONTENT |
      | feature-3 | local, origin | feature-3 commit | file      | content 3    |
    And Git Town setting "sync-feature-strategy" is "rebase"
    And the current branch is "feature-2"
    When I run "git-town delete"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                    |
      | feature-2 | git fetch --prune --tags   |
      |           | git push origin :feature-2 |
      |           | git checkout main          |
      | main      | git branch -D feature-2    |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES                   |
      | local, origin | main, feature-1, feature-3 |
    And this lineage exists now
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-3 | feature-1 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git branch feature-2 {{ sha 'feature-2 commit' }} |
      |        | git push -u origin feature-2                      |
      |        | git checkout feature-2                            |
    And the current branch is still "feature-2"
    And the branches are now
      | REPOSITORY    | BRANCHES                              |
      | local, origin | main, feature-1, feature-2, feature-3 |
    And the initial lineage exists now
