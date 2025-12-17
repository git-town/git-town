Feature: sync only branches whose remote is gone

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | main   | local, origin |
      | gamma | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "alpha"
    And origin ships the "beta" branch using the "squash-merge" ship-strategy

  @this
  Scenario: with "merge" feature sync strategy
    When I run "git-town sync --gone"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | alpha  | git fetch --prune --tags                          |
      |        | git checkout main                                 |
      | main   | git -c rebase.updateRefs=false rebase origin/main |
      |        | git checkout alpha                                |
      | alpha  | git merge --no-edit --ff main                     |
      |        | git push                                          |
      |        | git branch -D beta                                |
      |        | git checkout gamma                                |
      | gamma  | git merge --no-edit --ff main                     |
      |        | git push                                          |
      |        | git checkout alpha                                |
      | alpha  | git push --tags                                   |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                        |
      | main   | local, origin | beta commit                    |
      | alpha  | local, origin | alpha commit                   |
      |        |               | Merge branch 'main' into alpha |
      | gamma  | local, origin | gamma commit                   |
      |        |               | Merge branch 'main' into gamma |

  Scenario: with "rebase" feature sync strategy
    Given Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town sync --all"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                                      |
      | alpha      | git fetch --prune --tags                                                     |
      |            | git checkout main                                                            |
      | main       | git -c rebase.updateRefs=false rebase origin/main                            |
      |            | git checkout alpha                                                           |
      | alpha      | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |            | git push --force-with-lease --force-if-includes                              |
      |            | git checkout beta                                                            |
      | beta       | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
      |            | git push --force-with-lease --force-if-includes                              |
      |            | git checkout observed                                                        |
      | observed   | git -c rebase.updateRefs=false rebase origin/observed                        |
      |            | git checkout production                                                      |
      | production | git -c rebase.updateRefs=false rebase origin/production                      |
      |            | git push                                                                     |
      |            | git checkout qa                                                              |
      | qa         | git -c rebase.updateRefs=false rebase origin/qa                              |
      |            | git push                                                                     |
      |            | git checkout alpha                                                           |
      | alpha      | git push --tags                                                              |
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE                  |
      | main       | local, origin | main commit              |
      | alpha      | local, origin | alpha commit             |
      | beta       | local, origin | beta commit              |
      | parked     | local         | local parked commit      |
      |            | origin        | origin parked commit     |
      | observed   | local, origin | origin observed commit   |
      |            | local         | local observed commit    |
      | production | local, origin | origin production commit |
      |            |               | local production commit  |
      | qa         | local, origin | qa origin commit         |
      |            |               | qa local commit          |
