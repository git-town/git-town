Feature: sync the current omni feature branch using the "compress" sync-feature strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And Git Town setting "sync-feature-strategy" is "compress"
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               | FILE NAME           | FILE CONTENT           |
      | main    | local    | local main commit     | main local file     | main local content     |
      |         | origin   | origin main commit    | main origin file    | main origin content    |
      | feature | local    | local feature commit  | feature local file  | feature local content  |
      |         | origin   | origin feature commit | feature origin file | feature origin content |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git reset --soft main                   |
      |         | git commit -m "local feature commit"    |
      |         | git push --force-with-lease             |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE              |
      | main    | local, origin | origin main commit   |
      |         |               | local main commit    |
      | feature | local, origin | local feature commit |
    And these committed files exist now
      | BRANCH  | NAME                | CONTENT                |
      | main    | main local file     | main local content     |
      |         | main origin file    | main origin content    |
      | feature | feature local file  | feature local content  |
      |         | feature origin file | feature origin content |
      |         | main local file     | main local content     |
      |         | main origin file    | main origin content    |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
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
