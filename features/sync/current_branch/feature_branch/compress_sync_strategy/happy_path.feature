Feature: one person making a series of commits and syncs in between
  This scenario demonstrates that the "compress" strategy works
  as long as only one person contributes to a branch
  even if they change already committed content.

  Scenario:
    Given a Git repo with origin
    And the committed configuration file:
      """
      [branches]
      main = "main"
      perennials = []

      [sync]
      feature-strategy = "compress"
      """
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    #
    # I make the first change and sync
    Given I add this commit to the current branch:
      | MESSAGE     | FILE NAME | FILE CONTENT |
      | the feature | file      | content 1    |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | feature | local, origin | the feature | file      | content 1    |
    And all branches are now synchronized
    #
    # I make another change and sync
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME | FILE CONTENT |
      | my second commit | file      | content 2    |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | feature | local, origin | the feature | file      | content 2    |
    And all branches are now synchronized
    #
    # I make a third change and sync
    Given I add this commit to the current branch:
      | MESSAGE          | FILE NAME | FILE CONTENT |
      | my second commit | file      | content 3    |
    And wait 1 second to ensure new Git timestamps
    When I run "git-town sync"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/feature |
      |         | git reset --soft main                   |
      |         | git commit -m "the feature"             |
      |         | git push --force-with-lease             |
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | feature | local, origin | the feature | file      | content 3    |
    And all branches are now synchronized
