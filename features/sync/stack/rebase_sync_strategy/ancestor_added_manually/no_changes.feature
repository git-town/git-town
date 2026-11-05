@this
Feature: the parent branch was added manually

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-2 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | branch-2 | local, origin | commit 1 | file_1    | content 1    |
      | branch-2 | local, origin | commit 2 | file_2    | content 2    |
    And the commits
      | BRANCH | LOCATION | MESSAGE     | FILE NAME | FILE CONTENT |
      | main   | origin   | main commit | main_file | main         |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "branch-2"
    And I ran "git-town sync"
    And I ran "git branch branch-1 branch-2~1"
    And I ran "git config git-town-branch.branch-1.parent main"
    And I ran "git config git-town-branch.branch-2.parent branch-1"

  Scenario: result
    When I run "git-town sync --stack"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                                 |
      | branch-2 | git fetch --prune --tags                                                                |
      |          | git checkout main                                                                       |
      | main     | git -c rebase.updateRefs=false rebase origin/main                                       |
      |          | git checkout branch-1                                                                   |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto main {{ sha-before-run 'initial commit' }} |
      |          | git push -u origin branch-1                                                             |
      |          | git checkout branch-2                                                                   |
      | branch-2 | git -c rebase.updateRefs=false rebase branch-1                                          |
      |          | git push --force-with-lease --force-if-includes                                         |
    And no rebase is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And the initial commits exist now
