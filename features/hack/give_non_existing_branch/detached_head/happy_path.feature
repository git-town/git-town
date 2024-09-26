Feature: on a feature branch with a clean workspace

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the current branch is "existing"
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And I ran "git checkout HEAD^"
    When I run "git-town hack new"

  Scenario: result
    Then it runs the commands
      | BRANCH                     | COMMAND                  |
      | {{ sha 'initial commit' }} | git fetch --prune --tags |
      |                            | git checkout main        |
      | main                       | git rebase origin/main   |
      |                            | git checkout -b new      |
    And the current branch is now "new"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | main     | local, origin | main commit     |
      | existing | local         | existing commit |
      | new      | local         | main commit     |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH                     | COMMAND                                     |
      | new                        | git checkout main                           |
      | main                       | git reset --hard {{ sha 'initial commit' }} |
      |                            | git checkout {{ sha 'initial commit' }}     |
      | {{ sha 'initial commit' }} | git branch -D new                           |
    And the currently checked out commit is "initial commit"
