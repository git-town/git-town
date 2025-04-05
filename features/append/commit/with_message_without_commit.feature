Feature: message without commit

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
    And an uncommitted file with name "new_file" and content "new content"
    And I ran "git add new_file"
    When I run "git-town append new -m unrelated"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                 |
      | existing | git checkout -b new     |
      | new      | git commit -m unrelated |
      |          | git checkout existing   |
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
      | new      | local    | unrelated       |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | new      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND           |
      | existing | git branch -D new |
    And the initial commits exist now
    And the initial branches and lineage exist now
