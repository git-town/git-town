Feature: shipped changes conflict with multiple existing feature branches

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
      | gamma | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME        | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | conflicting_file | alpha content |
      | beta   | local, origin | beta commit  | conflicting_file | beta content  |
      | gamma  | local, origin | gamma commit | conflicting_file | gamma content |
    And origin ships the "beta" branch using the "squash-merge" ship-strategy
    And the current branch is "main"
    When I run "git-town sync --all"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git fetch --prune --tags                          |
      |        | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout alpha                                |
      | alpha  | git merge --no-edit --ff main                     |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a merge is now in progress
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                       |
      | alpha  | git commit --no-edit          |
      |        | git push                      |
      |        | git branch -D beta            |
      |        | git checkout gamma            |
      | gamma  | git merge --no-edit --ff main |
    And Git Town prints the error:
      """
      deleted branch beta
      """
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | gamma  | git commit --no-edit |
      |        | git push             |
      |        | git checkout main    |
      | main   | git push --tags      |
    And all branches are now synchronized
    And no merge is now in progress
    And these committed files exist now
      | BRANCH | NAME             | CONTENT          |
      | main   | conflicting_file | beta content     |
      | alpha  | conflicting_file | resolved content |
      | gamma  | conflicting_file | resolved content |
