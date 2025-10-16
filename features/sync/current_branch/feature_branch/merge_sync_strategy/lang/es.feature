Feature: sync the current branch in Spanish

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      | feature | local    | local feature commit |
    And the current branch is "feature"
    When I run "git-town sync" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git fetch --prune --tags                          |
      |         | git checkout main                                 |
      | main    | git -c rebase.updateRefs=false rebase origin/main |
      |         | git push                                          |
      |         | git checkout feature                              |
      | feature | git merge --no-edit --ff main                     |
      |         | git push -u origin feature                        |
    And Git Town prints:
      """
      Cambiado a rama 'feature'
      """
    And the branches are now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | local main commit                |
      | feature | local, origin | local feature commit             |
      |         |               | Merge branch 'main' into feature |

  Scenario: undo
    When I run "git-town undo" with these environment variables
      | LANG | es_ES.UTF-8 |
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git reset --hard {{ sha 'local feature commit' }} |
      |         | git push origin :feature                          |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE              |
      | main    | local, origin | local main commit    |
      | feature | local         | local feature commit |
