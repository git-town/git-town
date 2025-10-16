Feature: commit without message

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And the current branch is "existing"
    And an uncommitted file "new_file" with content "new content"
    And I ran "git add new_file"
    When I run "git-town hack new --commit" and enter "unrelated idea" for the commit message

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | existing | git checkout -b new main |
      | new      | git commit               |
      |          | git checkout existing    |
    And this lineage exists now
      """
      main
        existing
        new
      """
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
      | new      | local    | unrelated idea  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND           |
      | existing | git branch -D new |
    And the initial branches and lineage exist now
    And the initial commits exist now
