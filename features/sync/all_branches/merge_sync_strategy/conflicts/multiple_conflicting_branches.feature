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

  Scenario: skipping all conflicts
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH | COMMAND                       |
      | alpha  | git merge --abort             |
      |        | git checkout beta             |
      | beta   | git merge --no-edit --ff main |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in file
      """
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
    When I run "git-town skip"
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | gamma  | git merge --abort |
      |        | git checkout main |
      | main   | git push --tags   |
    And no merge is now in progress
