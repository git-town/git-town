Feature: using the "compress" strategy, sync a branch with unshipped local changes whose tracking branch was deleted

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | shipped | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE          |
      | shipped | local, origin | shipped commit   |
      |         | local         | unshipped commit |
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And origin ships the "shipped" branch using the "squash-merge" ship-strategy
    And the current branch is "shipped"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | shipped | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git checkout shipped                              |
      | shipped | git merge --no-edit --ff main                     |
    And Git Town prints:
      """
      Branch "shipped" was deleted at the remote but the local branch contains unshipped changes.
      """
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                               |
      | shipped | git checkout main                                     |
      | main    | git reset --hard {{ sha 'initial commit' }}           |
      |         | git checkout shipped                                  |
      | shipped | git reset --hard {{ sha-initial 'unshipped commit' }} |
    And the initial branches and lineage exist now
    And the initial commits exist now
