Feature: with upstream repo

  Background:
    Given a Git repo with origin
    And an upstream repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | main    | upstream | upstream commit |
      | feature | local    | local commit    |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                      |
      | feature | git fetch --prune --tags                                                     |
      |         | git checkout main                                                            |
      | main    | git fetch upstream main                                                      |
      |         | git -c rebase.updateRefs=false rebase upstream/main                          |
      |         | git push                                                                     |
      |         | git checkout feature                                                         |
      | feature | git push --force-with-lease --force-if-includes                              |
      |         | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |         | git push --force-with-lease --force-if-includes                              |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, origin, upstream | upstream commit |
      | feature | local, origin           | local commit    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                               |
      | feature | git reset --hard {{ sha-initial 'local commit' }}                     |
      |         | git push --force-with-lease origin {{ sha 'initial commit' }}:feature |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, origin, upstream | upstream commit |
      | feature | local                   | local commit    |
