Feature: prepend a branch to a feature branch in detached mode with a new commit on the tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | branch | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | branch | local, origin | local commit  |
      |        | origin        | origin commit |
    And the current branch is "branch"
    When I run "git-town prepend parent --detached"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                |
      | branch | git fetch --prune --tags               |
      |        | git merge --no-edit --ff origin/branch |
      |        | git checkout -b parent main            |
    And this lineage exists now
      """
      main
        parent
          branch
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE       |
      | branch | local, origin | local commit  |
      |        |               | origin commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                   |
      | parent | git checkout branch                       |
      | branch | git reset --hard {{ sha 'local commit' }} |
      |        | git branch -D parent                      |
    And the initial lineage exists now
    And the initial commits exist now
