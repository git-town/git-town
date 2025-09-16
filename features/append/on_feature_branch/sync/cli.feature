Feature: disable syncing via CLI

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE                |
      | branch-1 | local    | local branch-1 commit  |
      | branch-1 | origin   | origin branch-1 commit |
    And the current branch is "branch-1"
    When I run "git-town append branch-2 --no-sync"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                  |
      | branch-1 | git fetch --prune --tags                 |
      |          | git merge --no-edit --ff origin/branch-1 |
      |          | git push                                 |
      |          | git checkout -b branch-2                 |
    And the initial commits exist now
    And this lineage exists now
      """
      main
        existing
          new
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND               |
      | new      | git checkout existing |
      | existing | git branch -D new     |
    And the initial commits exist now
    And the initial lineage exists now
