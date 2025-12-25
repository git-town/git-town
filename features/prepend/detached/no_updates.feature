Feature: prepend a branch to a feature branch in detached mode when there are no new commits

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | branch | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | branch | local, origin | branch commit |
    And the current branch is "branch"
    When I run "git-town prepend parent --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                     |
      | branch | git fetch --prune --tags    |
      |        | git checkout -b parent main |
    And this lineage exists now
      """
      main
        parent
          branch
      """
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout branch  |
      | branch | git branch -D parent |
    And the initial lineage exists now
    And the initial commits exist now
