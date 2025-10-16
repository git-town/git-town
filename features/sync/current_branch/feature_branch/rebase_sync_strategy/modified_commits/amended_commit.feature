Feature: rebase a branch that contains amended commits

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-1 | feature | main      | local, origin |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local, origin | commit 1a | file_1    | one          |
      | feature-2 | local, origin | commit 2  | file_2    | two          |
    And the current branch is "feature-2"
    And I ran "git-town sync"
    And I amend this commit
      | BRANCH    | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local    | commit 1b | file_1    | another one  |
    And the current branch is "feature-2"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                              |
      | feature-2 | git fetch --prune --tags                                                             |
      |           | git checkout feature-1                                                               |
      | feature-1 | git push --force-with-lease --force-if-includes                                      |
      |           | git checkout feature-2                                                               |
      | feature-2 | git -c rebase.updateRefs=false rebase --onto feature-1 {{ sha-initial 'commit 1a' }} |
      |           | git push --force-with-lease --force-if-includes                                      |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local, origin | commit 1b | file_1    | another one  |
      | feature-2 | local, origin | commit 2  | file_2    | two          |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                              |
      | feature-2 | git reset --hard {{ sha 'commit 2' }}                                                |
      |           | git push --force-with-lease --force-if-includes                                      |
      |           | git push --force-with-lease origin {{ sha-in-origin-initial 'commit 1a' }}:feature-1 |
    And the branches are now
      | REPOSITORY    | BRANCHES                   |
      | local, origin | main, feature-1, feature-2 |
    And the initial commits exist now
