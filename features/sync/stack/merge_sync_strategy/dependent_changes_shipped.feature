Feature: shipped the head branch of a synced stack with dependent changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | file      | alpha content |
    And the branches
      | NAME | TYPE    | PARENT | LOCATIONS     |
      | beta | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | beta   | local, origin | beta commit | file      | beta content |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | file      | alpha content |
      | beta   | local, origin | beta commit  | file      | beta content  |
    And the current branch is "beta"
    And origin ships the "alpha" branch
    When I run "git-town sync"

  @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                 |
      | beta   | git fetch --prune --tags                |
      |        | git checkout main                       |
      | main   | git rebase origin/main --no-update-refs |
      |        | git branch -D alpha                     |
      |        | git checkout beta                       |
      | beta   | git merge --no-edit --ff origin/beta    |
      |        | git merge --no-edit --ff main           |
    # TODO: resolve this phantom merge conflict automatically. "file" has this content:
    #
    # <<<<<<< HEAD
    # beta content
    # =======
    # alpha content
    # >>>>>>> main
    #
    # It should choose "beta content" here.
    # Branch "beta" changes "alpha content" to "beta content".
    # This merge conflict is asking us to verify this again. It should not do that. It should know that the change from "alpha content" to "beta content" is legit.
    # Branch "alpha" has "alpha content" and branch "main" also has "alpha content" --> no unrelated changes, it's okay to use the version on "beta" here.
    #
    # Both these commands display the file content on various branches, even in the middle of a merge conflict:
    # 1. git show alpha:file
    # 2. git show 123456:file   (123456 is the SHA of a commit, for example the SHA that branch "alpha" points to)
    #
    # Another possible way to make this easier is to delete branch "alpha" at the end, after syncing all branches. This way, branch "alpha" is still around for checking the file content on it.
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a merge is now in progress
    And inspect the repo

  Scenario: resolve and continue
    When I resolve the conflict in "file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND              |
      | beta   | git commit --no-edit |
      |        | git push             |
    And the current branch is still "beta"
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                       | FILE NAME | FILE CONTENT     |
      | main   | local, origin | alpha commit                  | file      | alpha content    |
      | beta   | local, origin | alpha commit                  | file      | alpha content    |
      |        |               | beta commit                   | file      | beta content     |
      |        |               | Merge branch 'main' into beta | file      | resolved content |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                              |
      | beta   | git reset --hard {{ sha-before-run 'beta commit' }}  |
      |        | git push --force-with-lease --force-if-includes      |
      |        | git checkout main                                    |
      | main   | git reset --hard {{ sha 'initial commit' }}          |
      |        | git branch alpha {{ sha-before-run 'alpha commit' }} |
      |        | git checkout beta                                    |
    And the current branch is still "beta"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME | FILE CONTENT  |
      | main   | origin        | alpha commit | file      | alpha content |
      | alpha  | local         | alpha commit | file      | alpha content |
      | beta   | local, origin | beta commit  | file      | beta content  |
      |        | origin        | alpha commit | file      | alpha content |
    And the initial branches and lineage exist now
