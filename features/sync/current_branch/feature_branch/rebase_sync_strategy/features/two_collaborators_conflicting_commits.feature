Feature: collaborative feature branch syncing

  @debug @this
  Scenario: I and my coworker work together on a branch
    Given the committed configuration file:
      """
      [sync-strategy]
      feature-branches = "rebase"

      [branches]
      main = "main"
      perennials = []
      """
    And the current branch is a feature branch "feature"
    And a coworker clones the repository
    And the coworker fetches updates
    And the coworker is on the "feature" branch
    And the coworker sets the parent branch of "feature" as "main"
    And the commits
      | BRANCH  | LOCATION | MESSAGE   | FILE NAME | FILE CONTENT  |
      | feature | local    | my commit | file.txt  | my content 01 |
    When I run "git town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main                          |
      |         | git checkout feature                            |
      | feature | git rebase main                                 |
      |         | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE     |
      | main    | local, coworker, origin | config file |
      | feature | local, coworker, origin | config file |
      |         | local, origin           | my commit   |
    And all branches are now synchronized

    # coworker makes a conflicting local commit concurrently with me
    Given the coworker adds this commit:
      | MESSAGE         | FILE NAME | FILE CONTENT        |
      | coworker commit | file.txt  | coworker content 01 |
    When the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main                          |
      |         | git checkout feature                            |
      | feature | git rebase main                                 |
      |         | git push --force-with-lease --force-if-includes |
      |         | git rebase origin/feature                       |
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    When the coworker resolves the conflict in "file.txt" with "my and coworker content 01"
    And the coworker runs "git town continue" and closes the editor
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git rebase --continue                           |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, coworker, origin | config file     |
      | feature | local, coworker, origin | config file     |
      |         |                         | my commit       |
      |         | coworker, origin        | coworker commit |

    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                   |
      | feature | git fetch --prune --tags  |
      |         | git checkout main         |
      | main    | git rebase origin/main    |
      |         | git checkout feature      |
      | feature | git rebase origin/feature |
      |         | git rebase main           |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | feature | local, coworker, origin | my commit       |
      |         |                         | coworker commit |
