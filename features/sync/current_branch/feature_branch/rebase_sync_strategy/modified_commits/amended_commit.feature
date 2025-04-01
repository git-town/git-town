Feature: rebase a branch that contains amended commits

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE    | PARENT | LOCATIONS |
      | feature-1 | feature | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local    | commit 1a | file_1    | one          |
    And the branch
      | NAME      | TYPE    | PARENT    | LOCATIONS |
      | feature-2 | feature | feature-1 | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature-2 | local    | commit 2 | file_2    | two          |
    And the current branch is "feature-1"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And I ran "git town sync"
    And I amend this commit
      | BRANCH    | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT |
      | feature-1 | local    | commit 1b | file_1    | another one  |
    When I run "git-town sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main --no-update-refs         |
      |         | git push                                        |
      |         | git checkout feature                            |
      | feature | git rebase main --no-update-refs                |
      |         | git push --force-with-lease --force-if-includes |
      |         | git rebase origin/feature --no-update-refs      |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | origin main commit    |
      |         |               | local main commit     |
      | feature | local, origin | origin feature commit |
      |         |               | origin main commit    |
      |         |               | local main commit     |
      |         |               | local feature commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                           |
      | feature | git reset --hard {{ sha-before-run 'local feature commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin-before-run 'origin feature commit' }}:feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | origin main commit    |
      |         |               | local main commit     |
      | feature | local         | local feature commit  |
      |         | origin        | origin feature commit |
    And the initial branches and lineage exist now
