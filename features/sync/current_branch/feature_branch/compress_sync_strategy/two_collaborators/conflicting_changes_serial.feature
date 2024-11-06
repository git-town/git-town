Feature: two people make alternating conflicting changes to the same branch using the "compress" strategy

  This feature spec demonstrates a limitation of the "compress" sync strategy:
  If two people make conflicting changes to the same branch,
  they'll have to re-resolve merge conflicts
  even if they coordinate to avoid concurrent updates
  and run "git town sync" before and after they make changes.

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
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And a coworker clones the repository
    And the coworker fetches updates
    And the coworker is on the "feature" branch
    And the coworker sets the parent branch of "feature" as "main"

    # I add the first commit to the "feature" branch
    Given I add this commit to the current branch:
      | MESSAGE     | FILE NAME        | FILE CONTENT |
      | the feature | conflicting_file | my content 1 |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME        | FILE CONTENT |
      | feature | local, origin | the feature | conflicting_file | my content 1 |
    And all branches are now synchronized

    # my coworker syncs and adds a commit to the branch
    Given wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE     | FILE NAME        | FILE CONTENT |
      | feature | local, coworker, origin | the feature | conflicting_file | my content 1 |
    And all branches are now synchronized
    And the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME        | FILE CONTENT                        |
      | coworker first commit | conflicting_file | my content 1 and coworker content 1 |
    And wait 1 second to ensure new Git timestamps
    And the coworker runs "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME        | FILE CONTENT                        |
      | feature | local            | the feature | conflicting_file | my content 1                        |
      |         | coworker, origin | the feature | conflicting_file | my content 1 and coworker content 1 |
    And all branches are now synchronized

    # I sync, make another change, and sync again
    Given wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I resolve the conflict in "conflicting_file" with "my content 1 and coworker content 1"
    And I run "git town continue" and close the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git commit --no-edit        |
      |         | git reset --soft main       |
      |         | git commit -m "the feature" |
      |         | git push --force-with-lease |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME        | FILE CONTENT                        |
      | feature | local, origin | the feature | conflicting_file | my content 1 and coworker content 1 |
      |         | coworker      | the feature | conflicting_file | my content 1 and coworker content 1 |
    And all branches are now synchronized
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME        | FILE CONTENT                        |
      | my second commit | conflicting_file | my content 2 and coworker content 1 |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME        | FILE CONTENT                        |
      | feature | local, origin | the feature | conflicting_file | my content 2 and coworker content 1 |
      |         | coworker      | the feature | conflicting_file | my content 1 and coworker content 1 |
    And all branches are now synchronized

    # the coworker syncs, makes another change, and syncs again
    Given wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When the coworker resolves the conflict in "conflicting_file" with "my content 2 and coworker content 1"
    And the coworker runs "git town continue" and closes the editor
    Then Git Town runs the commands
      | BRANCH  | COMMAND                     |
      | feature | git commit --no-edit        |
      |         | git reset --soft main       |
      |         | git commit -m "the feature" |
      |         | git push --force-with-lease |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE     |
      | feature | local            | the feature |
      |         | coworker, origin | the feature |
    And all branches are now synchronized
    Given the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME        | FILE CONTENT                        |
      | coworker first commit | conflicting_file | my content 2 and coworker content 2 |
    And wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME        | FILE CONTENT                        |
      | feature | local            | the feature | conflicting_file | my content 2 and coworker content 1 |
      |         | coworker, origin | the feature | conflicting_file | my content 2 and coworker content 2 |
    And all branches are now synchronized
