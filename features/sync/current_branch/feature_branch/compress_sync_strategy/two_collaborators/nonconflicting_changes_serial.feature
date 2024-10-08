Feature: two people make alternating non-conflicting changes to the same branch using the "compress" strategy

  This scenario demonstrates a limitation of the "compress" sync strategy:
  If two people collaborate on the same branch,
  they will run into many merge conflicts
  because each of their branches contains a single commit that introduces all the changes
  and Git doesn't know which of the branches is correct.

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
      | MESSAGE     | FILE NAME | FILE CONTENT | FILE NAME | FILE CONTENT                         |
      | the feature | my_file_1 | my content 1 | file      | my content 1 \n\n coworker content 0 |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                         |
      | feature | local, origin | the feature | file      | my content 1 \n\n coworker content 0 |
    And all branches are now synchronized

    # my coworker syncs and adds a commit to the branch
    Given wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE     | FILE NAME | FILE CONTENT                         |
      | feature | local, coworker, origin | the feature | file      | my content 1 \n\n coworker content 0 |
    And all branches are now synchronized
    And the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME | FILE CONTENT                         |
      | coworker first commit | file      | my content 1 \n\n coworker content 1 |
    And wait 1 second to ensure new Git timestamps
    And the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME | FILE CONTENT                         |
      | feature | local            | the feature | file      | my content 1 \n\n coworker content 0 |
      |         | coworker, origin | the feature | file      | my content 1 \n\n coworker content 1 |
    And all branches are now synchronized

    # I sync, make another change, and sync again
    Given wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    When I resolve the conflict in "file" with "my content 1 \n\n coworker content 1"
    And I run "git town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git commit --no-edit          |
      |         | git merge --no-edit --ff main |
      |         | git reset --soft main         |
      |         | git commit -m "the feature"   |
      |         | git push --force-with-lease   |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                         |
      | feature | local, origin | the feature | file      | my content 1 \n\n coworker content 1 |
      |         | coworker      | the feature | file      | my content 1 \n\n coworker content 1 |
    And all branches are now synchronized
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME | FILE CONTENT                         |
      | my second commit | file      | my content 2 \n\n coworker content 1 |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT                         |
      | feature | local, origin | the feature | file      | my content 2 \n\n coworker content 1 |
      |         | coworker      | the feature | file      | my content 1 \n\n coworker content 1 |
    And all branches are now synchronized

    # the coworker syncs, makes another change, and syncs again
    Given wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    When the coworker resolves the conflict in "file" with "my content 2 \n\n coworker content 1"
    And the coworker runs "git town continue" and closes the editor
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git commit --no-edit          |
      |         | git merge --no-edit --ff main |
      |         | git reset --soft main         |
      |         | git commit -m "the feature"   |
      |         | git push --force-with-lease   |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE     |
      | feature | local            | the feature |
      |         | coworker, origin | the feature |
    And all branches are now synchronized
    Given the coworker adds this commit to their current branch:
      | MESSAGE               | FILE NAME | FILE CONTENT                         |
      | coworker first commit | file      | my content 2 \n\n coworker content 2 |
    And wait 1 second to ensure new Git timestamps
    When the coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME | FILE CONTENT                         |
      | feature | local            | the feature | file      | my content 2 \n\n coworker content 1 |
      |         | coworker, origin | the feature | file      | my content 2 \n\n coworker content 2 |
    And all branches are now synchronized
