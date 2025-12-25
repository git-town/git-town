Feature: shipped branch with multiple descendents

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME      | FILE CONTENT      |
      | feature-1 | local, origin | feature-1 commit | feature-1-file | feature 1 content |
    And the branches
      | NAME       | TYPE    | PARENT    | LOCATIONS     |
      | feature-1a | feature | feature-1 | local, origin |
    And the commits
      | BRANCH     | LOCATION      | MESSAGE           | FILE NAME       | FILE CONTENT       |
      | feature-1a | local, origin | feature-1a commit | feature-1a-file | feature 1a content |
    And the branches
      | NAME       | TYPE    | PARENT    | LOCATIONS     |
      | feature-1b | feature | feature-1 | local, origin |
    And the commits
      | BRANCH     | LOCATION      | MESSAGE           | FILE NAME       | FILE CONTENT       |
      | feature-1b | local, origin | feature-1b commit | feature-1b-file | feature 1b content |
    And Git setting "git-town.sync-feature-strategy" is "merge"
    And origin ships the "feature-1" branch using the "squash-merge" ship-strategy
    And the current branch is "feature-1"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                           |
      | feature-1  | git fetch --prune --tags                          |
      |            | git checkout main                                 |
      | main       | git -c rebase.updateRefs=false rebase origin/main |
      |            | git branch -D feature-1                           |
      |            | git checkout feature-1a                           |
      | feature-1a | git merge --no-edit --ff main                     |
      |            | git push                                          |
      |            | git checkout feature-1b                           |
      | feature-1b | git merge --no-edit --ff main                     |
      |            | git push                                          |
      |            | git push --tags                                   |
    And Git Town prints:
      """
      deleted branch "feature-1"
      """
    And Git Town prints:
      """
      branch "feature-1a" is now a child of "main"
      """
    And Git Town prints:
      """
      branch "feature-1b" is now a child of "main"
      """
    And this lineage exists now
      """
      main
        feature-1a
        feature-1b
      """
    And the branches are now
      | REPOSITORY    | BRANCHES                     |
      | local, origin | main, feature-1a, feature-1b |
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                             |
      | main       | local, origin | feature-1 commit                    |
      | feature-1a | local, origin | feature-1 commit                    |
      |            |               | feature-1a commit                   |
      |            |               | Merge branch 'main' into feature-1a |
      | feature-1b | local, origin | feature-1 commit                    |
      |            |               | feature-1b commit                   |
      |            |               | Merge branch 'main' into feature-1b |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                           |
      | feature-1b | git checkout feature-1a                           |
      | feature-1a | git reset --hard {{ sha 'feature-1a commit' }}    |
      |            | git push --force-with-lease --force-if-includes   |
      |            | git checkout feature-1b                           |
      | feature-1b | git reset --hard {{ sha 'feature-1b commit' }}    |
      |            | git push --force-with-lease --force-if-includes   |
      |            | git checkout main                                 |
      | main       | git reset --hard {{ sha 'initial commit' }}       |
      |            | git branch feature-1 {{ sha 'feature-1 commit' }} |
      |            | git checkout feature-1                            |
    And the initial branches and lineage exist now
