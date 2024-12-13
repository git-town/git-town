Feature: two people with rebase strategy sync changes made by them

  Scenario: I and my coworker sync changes we both made to the same branch
    Given a Git repo with origin
    And the committed configuration file:
      """
      [branches]
      main = "main"
      perennials = []

      [sync]
      feature-strategy = "rebase"
      """
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And a coworker clones the repository
    And the current branch is "feature"
    And the coworker fetches updates
    And the coworker sets the parent branch of "feature" as "main"
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | feature | local    | my commit       |
      |         | coworker | coworker commit |
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main --no-update-refs         |
      |         | git checkout feature                            |
      | feature | git rebase main --no-update-refs                |
      |         | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         |
      | feature | local, origin | my commit       |
      |         | coworker      | coworker commit |
    And all branches are now synchronized

    Given the coworker is on the "feature" branch
    When the coworker runs "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main --no-update-refs         |
      |         | git checkout feature                            |
      | feature | git rebase main --no-update-refs                |
      |         | git push --force-with-lease --force-if-includes |
      |         | git rebase origin/feature --no-update-refs      |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | feature | local, coworker, origin | my commit       |
      |         | coworker, origin        | coworker commit |

    Given the current branch is "feature"
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main --no-update-refs         |
      |         | git checkout feature                            |
      | feature | git rebase main --no-update-refs                |
      |         | git push --force-with-lease --force-if-includes |
      |         | git rebase origin/feature --no-update-refs      |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | feature | local, coworker, origin | my commit       |
      |         |                         | coworker commit |
