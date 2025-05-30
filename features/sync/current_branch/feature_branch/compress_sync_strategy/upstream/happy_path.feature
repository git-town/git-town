Feature: "compress" sync with upstream repo

  Background:
    Given a Git repo with origin
    And an upstream repo
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME     | FILE CONTENT     |
      | main    | upstream | upstream commit | upstream_file | upstream content |
      | feature | local    | local commit    | local file    | local content    |
    And the current branch is "feature"
    And Git setting "git-town.sync-feature-strategy" is "compress"
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                             |
      | feature | git fetch --prune --tags                            |
      |         | git checkout main                                   |
      | main    | git fetch upstream main                             |
      |         | git -c rebase.updateRefs=false rebase upstream/main |
      |         | git push                                            |
      |         | git checkout feature                                |
      | feature | git merge --no-edit --ff main                       |
      |         | git merge --no-edit --ff origin/feature             |
      |         | git reset --soft main                               |
      |         | git commit -m "local commit"                        |
      |         | git push --force-with-lease                         |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         | FILE NAME     | FILE CONTENT     |
      | main    | local, origin, upstream | upstream commit | upstream_file | upstream content |
      | feature | local, origin           | local commit    | local file    | local content    |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                               |
      | feature | git reset --hard {{ sha-initial 'local commit' }}                     |
      |         | git push --force-with-lease origin {{ sha 'initial commit' }}:feature |
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, origin, upstream | upstream commit |
      | feature | local                   | local commit    |
    And the initial branches and lineage exist now
