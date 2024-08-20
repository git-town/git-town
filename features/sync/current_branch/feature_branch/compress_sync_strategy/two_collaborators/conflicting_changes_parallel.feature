Feature: two people using the "compress" strategy make concurrent conflicting changes to a branch

  This feature spec demonstrates what happens
  when two people make concurrent changes to the same branch
  and run "git sync" before and after they make changes.
  Running "git sync" so often surfaces merge conflicts early when they are still small and easy to resolve.
  A downside of the "compress" strategy is that both committers keep renaming the commit on the branch
  to the first commit that exists on their local branch, which is different in this example.

  Scenario:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [sync-strategy]
      feature-branches = "compress"

      [branches]
      main = "main"
      perennials = []
      """
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And a coworker clones the repository
    And the coworker fetches updates
    And the coworker is on the "feature" branch
    And the coworker sets the parent branch of "feature" as "main"

    # I make a commit and sync
    Given I add this commit to the current branch:
      | MESSAGE         | FILE NAME        | FILE CONTENT |
      | my first commit | conflicting_file | my content 1 |
    And wait 1 second to ensure new Git timestamps
    When I run "git town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git reset --soft main                   |
      |         | git commit -m "my first commit"         |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE         | FILE NAME        | FILE CONTENT |
      | feature | local, origin | my first commit | conflicting_file | my content 1 |
    And all branches are now synchronized

    # the coworker makes a conflicting local commit concurrently with me and syncs
    Given the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME        | FILE CONTENT       |
      | coworker first commit | conflicting_file | coworker content 1 |
    And wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When the coworker resolves the conflict in "conflicting_file" with "my content 1 and coworker content 1"
    And the coworker runs "git town continue" and closes the editor
    Then it runs the commands
      | BRANCH  | COMMAND                               |
      | feature | git commit --no-edit                  |
      |         | git merge --no-edit --ff main         |
      |         | git reset --soft main                 |
      |         | git commit -m "coworker first commit" |
      |         | git push --force-with-lease           |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE               | FILE NAME        | FILE CONTENT                        |
      | feature | local            | my first commit       | conflicting_file | my content 1                        |
      |         | coworker, origin | coworker first commit | conflicting_file | my content 1 and coworker content 1 |

    # I add another conflicting commit locally and then sync
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME        | FILE CONTENT |
      | my second commit | conflicting_file | my content 2 |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I resolve the conflict in "conflicting_file" with "my content 2 and coworker content 1"
    And I run "git town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git commit --no-edit            |
      |         | git merge --no-edit --ff main   |
      |         | git reset --soft main           |
      |         | git commit -m "my first commit" |
      |         | git push --force-with-lease     |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE               | FILE NAME        | FILE CONTENT                        |
      | feature | local, origin | my first commit       | conflicting_file | my content 2 and coworker content 1 |
      |         | coworker      | coworker first commit | conflicting_file | my content 1 and coworker content 1 |

    # the coworker makes another conflicting local commit concurrently with me and syncs
    Given the coworker adds this commit to their current branch:
      | MESSAGE                | FILE NAME        | FILE CONTENT                        |
      | coworker second commit | conflicting_file | my content 1 and coworker content 2 |
    And wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When the coworker resolves the conflict in "conflicting_file" with "my content 2 and coworker content 2"
    And the coworker runs "git town continue" and closes the editor
    Then it runs the commands
      | BRANCH  | COMMAND                               |
      | feature | git commit --no-edit                  |
      |         | git merge --no-edit --ff main         |
      |         | git reset --soft main                 |
      |         | git commit -m "coworker first commit" |
      |         | git push --force-with-lease           |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE               | FILE NAME        | FILE CONTENT                        |
      | feature | local            | my first commit       | conflicting_file | my content 2 and coworker content 1 |
      |         | coworker, origin | coworker first commit | conflicting_file | my content 2 and coworker content 2 |
