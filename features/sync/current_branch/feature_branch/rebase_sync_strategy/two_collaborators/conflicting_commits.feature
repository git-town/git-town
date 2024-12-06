Feature: two people using rebase make conflicting changes to a branch

  Scenario: I and my coworker work together on a branch
    Given a Git repo with origin
    And the committed configuration file:
      """
      [sync-strategy]
      feature-branches = "rebase"

      [branches]
      main = "main"
      perennials = []
      """
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And a coworker clones the repository
    And the coworker fetches updates
    And the coworker is on the "feature" branch
    And the coworker sets the parent branch of "feature" as "main"

    # I make a commit and sync
    Given I add this commit to the current branch:
      | MESSAGE         | FILE NAME | FILE CONTENT |
      | my first commit | file.txt  | my content   |
    When I run "git town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git checkout main                               |
      | main    | git rebase origin/main --no-update-refs         |
      |         | git checkout feature                            |
      | feature | git rebase main --no-update-refs                |
      |         | git push --force-with-lease --force-if-includes |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         | FILE NAME | FILE CONTENT |
      | feature | local, origin | my first commit | file.txt  | my content   |
    And all branches are now synchronized

    # coworker makes a conflicting local commit concurrently with me and then syncs
    Given the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME | FILE CONTENT     |
      | coworker first commit | file.txt  | coworker content |
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
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    When the coworker resolves the conflict in "file.txt" with "my and coworker content"
    And the coworker runs "git town continue" and closes the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git -c core.editor=true rebase --continue       |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE               | FILE NAME | FILE CONTENT            |
      | feature | local, coworker, origin | my first commit       | file.txt  | my content              |
      |         | coworker, origin        | coworker first commit | file.txt  | my and coworker content |

    # I add a conflicting commit locally and then sync
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME | FILE CONTENT   |
      | my second commit | file.txt  | my new content |
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
    And Git Town prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      """
    When I resolve the conflict in "file.txt" with "my new and coworker content"
    And I run "git town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git -c core.editor=true rebase --continue       |
      |         | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE               | FILE NAME | FILE CONTENT                |
      | feature | local, coworker, origin | my first commit       | file.txt  | my content                  |
      |         |                         | coworker first commit | file.txt  | my and coworker content     |
      |         | local, origin           | my second commit      | file.txt  | my new and coworker content |
