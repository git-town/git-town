Feature: detaching a local branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      | branch-1 | local    | commit 1b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-2 | feature | branch-1 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-2 | local    | commit 2a |
      | branch-2 | local    | commit 2b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-3 | feature | branch-2 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-3 | local    | commit 3a |
      | branch-3 | local    | commit 3b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-4 | feature | branch-3 | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-4 | local    | commit 4a |
      | branch-4 | local    | commit 4b |
    And the current branch is "branch-2"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                             |
      | branch-2 | git fetch --prune --tags            |
      |          | git rebase --onto main branch-1     |
      |          | git checkout branch-3               |
      | branch-3 | git rebase --onto branch-1 branch-2 |
      |          | git checkout branch-4               |
      | branch-4 | git rebase --onto branch-3 branch-2 |
      |          | git checkout branch-2               |
    And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      |          |          | commit 1b |
      | branch-2 | local    | commit 2a |
      |          |          | commit 2b |
      | branch-3 | local    | commit 3a |
      |          |          | commit 3b |
      | branch-4 | local    | commit 4a |
      |          |          | commit 4b |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | main     |
      | branch-3 | branch-1 |
      | branch-4 | branch-3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                |
      | branch-2 | git reset --hard {{ sha 'commit 2b' }} |
      |          | git checkout branch-3                  |
      | branch-3 | git reset --hard {{ sha 'commit 3b' }} |
      |          | git checkout branch-4                  |
      | branch-4 | git reset --hard {{ sha 'commit 4b' }} |
      |          | git checkout branch-2                  |
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial lineage exists now
