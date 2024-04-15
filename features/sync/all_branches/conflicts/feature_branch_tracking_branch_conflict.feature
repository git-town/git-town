Feature: handle merge conflicts between feature branches and their tracking branches

  Background:
    Given the feature branches "alpha", "beta", and "gamma"
    And the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main   | origin        | main commit        | main_file        | main content        |
      | alpha  | local, origin | alpha commit       | feature1_file    | alpha content       |
      | beta   | local         | local beta commit  | conflicting_file | local beta content  |
      |        | origin        | origin beta commit | conflicting_file | origin beta content |
      | gamma  | local, origin | gamma commit       | feature3_file    | gamma content       |
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                          |
      | main   | git fetch --prune --tags         |
      |        | git add -A                       |
      |        | git stash                        |
      |        | git rebase origin/main           |
      |        | git checkout alpha               |
      | alpha  | git merge --no-edit origin/alpha |
      |        | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout beta                |
      | beta   | git merge --no-edit origin/beta  |
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
    And the current branch is now "beta"
    And the uncommitted file is stashed
    And a merge is now in progress

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | beta   | git merge --abort                               |
      |        | git checkout alpha                              |
      | alpha  | git reset --hard {{ sha 'alpha commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git stash pop                                   |
    And the current branch is now "main"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: skip
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH | COMMAND                          |
      | beta   | git merge --abort                |
      |        | git checkout gamma               |
      | gamma  | git merge --no-edit origin/gamma |
      |        | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout main                |
      | main   | git push --tags                  |
      |        | git stash pop                    |
    And the current branch is now "main"
    And the uncommitted file still exists
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | main   | local, origin | main commit                    |
      | alpha  | local, origin | alpha commit                   |
      |        |               | main commit                    |
      |        |               | Merge branch 'main' into alpha |
      | beta   | local         | local beta commit              |
      |        | origin        | origin beta commit             |
      | gamma  | local, origin | gamma commit                   |
      |        |               | main commit                    |
      |        |               | Merge branch 'main' into gamma |
    And these committed files exist now
      | BRANCH | NAME             | CONTENT            |
      | main   | main_file        | main content       |
      | alpha  | feature1_file    | alpha content      |
      |        | main_file        | main content       |
      | beta   | conflicting_file | local beta content |
      | gamma  | feature3_file    | gamma content      |
      |        | main_file        | main content       |

  Scenario: continue with unresolved conflict
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And the current branch is still "beta"
    And the uncommitted file is stashed
    And a merge is now in progress

  Scenario: resolve and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH | COMMAND                          |
      | beta   | git commit --no-edit             |
      |        | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout gamma               |
      | gamma  | git merge --no-edit origin/gamma |
      |        | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout main                |
      | main   | git push --tags                  |
      |        | git stash pop                    |
    And all branches are now synchronized
    And the current branch is now "main"
    And the uncommitted file still exists
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
    Then it runs the commands
      | BRANCH | COMMAND                          |
      | beta   | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout gamma               |
      | gamma  | git merge --no-edit origin/gamma |
      |        | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout main                |
      | main   | git push --tags                  |
      |        | git stash pop                    |
