Feature: on a detached head with a clean workspace

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
    And I ran "git checkout HEAD^"
    When I run "git-town hack new"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH                     | COMMAND                                           |
      | {{ sha 'initial commit' }} | git fetch --prune --tags                          |
      |                            | git checkout main                                 |
      | main                       | git -c rebase.updateRefs=false rebase origin/main |
      |                            | git checkout -b new                               |
    And this lineage exists now
      """
      main
        existing
        new
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | main     | local, origin | main commit     |
      | existing | local         | existing commit |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH                     | COMMAND                                     |
      | new                        | git checkout main                           |
      | main                       | git reset --hard {{ sha 'initial commit' }} |
      |                            | git checkout {{ sha 'initial commit' }}     |
      | {{ sha 'initial commit' }} | git branch -D new                           |
