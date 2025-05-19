Feature: multiple conflicting branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
      | gamma | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME | FILE CONTENT        |
      | main   | origin        | main commit        | file      | main content        |
      | alpha  | local, origin | alpha commit       | file      | alpha content       |
      | beta   | local         | local beta commit  | file      | local beta content  |
      |        | origin        | origin beta commit | file      | origin beta content |
      | gamma  | local, origin | gamma commit       | file      | gamma content       |
    And the current branch is "main"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git fetch --prune --tags                          |
      |        | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout alpha                                |
      | alpha  | git merge --no-edit --ff main                     |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                     |
      | alpha  | git merge --abort                           |
      |        | git checkout main                           |
      | main   | git reset --hard {{ sha 'initial commit' }} |
    And the initial commits exist now
    And the initial branches and lineage exist now

  @this
  Scenario: skipping all conflicts
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH | COMMAND                       |
      | alpha  | git merge --abort             |
      |        | git checkout beta             |
      | beta   | git merge --no-edit --ff main |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, origin | main commit        |
      | alpha  | local, origin | alpha commit       |
      | beta   | local         | local beta commit  |
      |        | origin        | origin beta commit |
      | gamma  | local, origin | gamma commit       |
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH | COMMAND                       |
      | beta   | git merge --abort             |
      |        | git checkout gamma            |
      | gamma  | git merge --no-edit --ff main |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
    And a merge is now in progress
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE            |
      | main   | local, origin | main commit        |
      | alpha  | local, origin | alpha commit       |
      | beta   | local         | local beta commit  |
      |        | origin        | origin beta commit |
      | gamma  | local, origin | gamma commit       |

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                       |
      | beta   | git commit --no-edit          |
      |        | git push                      |
      |        | git checkout gamma            |
      | gamma  | git merge --no-edit --ff main |
      |        | git push                      |
      |        | git checkout main             |
      | main   | git push --tags               |
    And all branches are now synchronized
    And no merge is in progress
    And these committed files exist now
      | BRANCH | NAME             | CONTENT          |
      | main   | main_file        | main content     |
      | alpha  | feature1_file    | alpha content    |
      |        | main_file        | main content     |
      | beta   | conflicting_file | resolved content |
      |        | main_file        | main content     |
      | gamma  | feature3_file    | gamma content    |
      |        | main_file        | main content     |

  Scenario: resolve, commit, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then Git Town runs the commands
      | BRANCH | COMMAND                       |
      | beta   | git push                      |
      |        | git checkout gamma            |
      | gamma  | git merge --no-edit --ff main |
      |        | git push                      |
      |        | git checkout main             |
      | main   | git push --tags               |
