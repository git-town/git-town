Feature: syncing a stack that contains an observed branch

  Background:
    Given a Git repo with origin
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the branches
      | NAME     | TYPE   | LOCATIONS |
      | observed | (none) | origin    |
    And the current branch is "main"
    And I ran "git fetch"
    And I ran "git-town observe observed"
    And the commits
      | BRANCH   | LOCATION | MESSAGE    |
      | observed | origin   | new commit |
    When I run "git-town sync --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                               |
      | observed | git fetch --prune --tags                              |
      |          | git -c rebase.updateRefs=false rebase origin/observed |
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE    |
      | observed | local, origin | new commit |
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                     |
      | observed | git reset --hard {{ sha 'initial commit' }} |
    And the branches are now
      | REPOSITORY    | BRANCHES       |
      | local, origin | main, observed |
    And the initial commits exist now
