Feature: collaborative feature branch syncing

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

    # I make a commit and sync
    Given I add this commit:
      | MESSAGE         | FILE NAME | FILE CONTENT |
      | my first commit | file.txt  | my content   |
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
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, coworker, origin | config file     |
      | feature | local, coworker, origin | config file     |
      |         | local, origin           | my first commit |
    And all branches are now synchronized

    # coworker makes a conflicting local commit concurrently with me and then syncs
    Given the coworker adds this commit:
      | MESSAGE               | FILE NAME | FILE CONTENT     |
      | coworker first commit | file.txt  | coworker content |
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
    When the coworker resolves the conflict in "file.txt" with "my and coworker content"
    And the coworker runs "git town continue" and closes the editor
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git rebase --continue                           |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE               |
      | main    | local, coworker, origin | config file           |
      | feature | local, coworker, origin | config file           |
      |         |                         | my first commit       |
      |         | coworker, origin        | coworker first commit |

    # I add a conflicting commit locally and then sync
    Given I add this commit:
      | MESSAGE          | FILE NAME | FILE CONTENT |
      | my second commit | file.txt  | my content   |
    When I run "git-town sync"
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
    When I resolve the conflict in "file.txt" with "my and coworker content 02"
    And I run "git town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git rebase --continue                           |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE               |
      | main    | local, coworker, origin | config file           |
      | feature | local, coworker, origin | config file           |
      |         |                         | my first commit       |
      |         |                         | coworker first commit |
      |         | local, origin           | my second commit      |
