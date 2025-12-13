Feature: skip and park the current branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
      | gamma | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME        | FILE CONTENT  |
      | main   | origin        | main commit  | conflicting_file | main content  |
      | alpha  | local, origin | alpha commit | feature1_file    | alpha content |
      | beta   | local, origin | beta commit  | conflicting_file | beta content  |
      | gamma  | local, origin | gamma commit | feature2_file    | gamma content |
    And the current branch is "main"
    And I run "git-town sync --all"
    And Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git fetch --prune --tags                          |
      |        | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout alpha                                |
      | alpha  | git merge --no-edit --ff main                     |
      |        | git push                                          |
      |        | git checkout beta                                 |
      | beta   | git merge --no-edit --ff main                     |
    When I run "git-town skip --park"

  @this
  Scenario: skip with --park flag
    Then Git Town runs the commands
      | BRANCH | COMMAND                       |
      | beta   | git merge --abort             |
      |        | git checkout gamma            |
      | gamma  | git merge --no-edit --ff main |
      |        | git push                      |
      |        | git checkout main             |
      | main   | git push --tags               |
    And Git Town prints:
      """
      branch "beta" is now parked
      """
    And no merge is now in progress
    And branch "beta" now has type "parked"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | main   | local, origin | main commit                    |
      | alpha  | local, origin | alpha commit                   |
      |        |               | Merge branch 'main' into alpha |
      | beta   | local, origin | beta commit                    |
      | gamma  | local, origin | gamma commit                   |
      |        |               | Merge branch 'main' into gamma |
