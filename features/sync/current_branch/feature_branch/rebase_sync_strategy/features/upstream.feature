Feature: with upstream repo

  Background:
    Given Git Town setting "sync-feature-strategy" is "rebase"
    And an upstream repo
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | main    | upstream | upstream commit |
      | feature | local    | local commit    |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main                          |
      |         | git fetch upstream main                         |
      |         | git rebase upstream/main                        |
      |         | git push                                        |
      |         | git checkout feature                            |
      | feature | git rebase origin/feature                       |
      |         | git rebase main                                 |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, origin, upstream | upstream commit |
      | feature | local, origin           | upstream commit |
      |         |                         | local commit    |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                               |
      | feature | git reset --hard {{ sha-before-run 'local commit' }}                  |
      |         | git push --force-with-lease origin {{ sha 'initial commit' }}:feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, origin, upstream | upstream commit |
      | feature | local                   | local commit    |
    And the initial branches and lineage exist
