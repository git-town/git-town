Feature: already existing remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | old      | feature | main   | local, origin |
      | existing | feature | main   | origin        |
    And the current branch is "old"
    And an uncommitted file
    When I run "git-town prepend existing"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                       |
      | old      | git add -A                    |
      |          | git stash -m "Git Town WIP"   |
      |          | git checkout -b existing main |
      | existing | git stash pop                 |
      |          | git restore --staged .        |
    And this lineage exists now
      """
      main
        existing
          old
      """
    And the uncommitted file still exists
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | existing | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout old            |
      | old      | git branch -D existing      |
      |          | git stash pop               |
      |          | git restore --staged .      |
    And the initial lineage exists now
    And the uncommitted file still exists
    And the initial commits exist now
