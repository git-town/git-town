Feature: prune an empty branch that has a child branch using the rebase sync strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-1 | feature | main      | local, origin |
      | feature-2 | feature | feature-1 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME | FILE CONTENT |
      | feature-2 | local, origin | feature-2 commit | file      | content      |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "feature-1"
    When I run "git-town sync --prune"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                     |
      | feature-1 | git fetch --prune --tags                                    |
      |           | git checkout feature-2                                      |
      | feature-2 | git push origin :feature-1                                  |
      |           | git -c rebase.updateRefs=false rebase --onto main feature-1 |
      |           | git branch -D feature-1                                     |
    And this lineage exists now
      """
      main
        feature-2
      """
    And the branches are now
      | REPOSITORY    | BRANCHES        |
      | local, origin | main, feature-2 |
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | feature-2 | git branch feature-1 {{ sha 'initial commit' }} |
      |           | git push -u origin feature-1                    |
      |           | git checkout feature-1                          |
    And the initial branches and lineage exist now
    And the initial commits exist now
