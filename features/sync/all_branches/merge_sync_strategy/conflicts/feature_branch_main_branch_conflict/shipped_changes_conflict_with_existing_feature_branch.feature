Feature: shipped changes conflict with multiple existing feature branches

  Scenario:
    Given the feature branches "alpha", "beta", and "gamma"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME        | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | conflicting_file | alpha content |
      | beta   | local, origin | beta commit  | conflicting_file | beta content  |
      | gamma  | local, origin | gamma commit | conflicting_file | gamma content |
    And origin ships the "beta" branch
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town sync --all"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | main   | git fetch --prune --tags              |
      |        | git add -A                            |
      |        | git stash                             |
      |        | git rebase origin/main                |
      |        | git checkout alpha                    |
      | alpha  | git merge --no-edit --ff origin/alpha |
      |        | git merge --no-edit --ff main         |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git town skip".
      """
    And the current branch is now "alpha"
    And the uncommitted file is stashed
    And a merge is now in progress
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND                               |
      | alpha  | git commit --no-edit                  |
      |        | git push                              |
      |        | git checkout beta                     |
      | beta   | git merge --no-edit --ff main         |
      |        | git checkout main                     |
      | main   | git branch -D beta                    |
      |        | git checkout gamma                    |
      | gamma  | git merge --no-edit --ff origin/gamma |
      |        | git merge --no-edit --ff main         |
    And it prints something like:
      """
      deleted branch "beta"
      """
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And it prints the error:
      """
      To continue after having resolved conflicts, run "git town continue".
      To go back to where you started, run "git town undo".
      To continue by skipping the current branch, run "git town skip".
      """
    And the current branch is now "gamma"
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | gamma  | git commit --no-edit |
      |        | git push             |
      |        | git checkout main    |
      | main   | git push --tags      |
      |        | git stash pop        |
    And the current branch is now "main"
    And the uncommitted file still exists
    And all branches are now synchronized
    And no merge is in progress
    And these committed files exist now
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | beta content     |
      | alpha  | conflicting_file | resolved content |
      | gamma  | conflicting_file | resolved content |
