Feature: the branch was shipped manually on the local machine

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.unknown-branch-type" is "prototype"
    And origin deletes the "feature" branch
    And the current branch is "main"
    And I ran "git merge feature --squash"
    And I ran "git commit -m merged"
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                           |
      | main   | git fetch --prune --tags                          |
      |        | git -c rebase.updateRefs=false rebase origin/main |
      |        | git push                                          |
      |        | git branch -D feature                             |
      |        | git push --tags                                   |
    And no lineage exists now
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
    And this lineage exists now
      """
      main
        feature
      """
    And the branches are now
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | origin     | main          |
