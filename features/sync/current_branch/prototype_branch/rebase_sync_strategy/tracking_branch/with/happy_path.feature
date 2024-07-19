Feature: sync the current prototype branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the current branch is "prototype"
    And the commits
      | BRANCH    | LOCATION | MESSAGE            |
      | main      | local    | main local commit  |
      |           | local    | main origin commit |
      | prototype | local    | local commit       |
      |           | origin   | origin commit      |
    And Git Town setting "sync-feature-strategy" is "rebase"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                     |
      | prototype | git fetch --prune --tags    |
      |           | git checkout main           |
      | main      | git rebase origin/main      |
      |           | git push                    |
      |           | git checkout prototype      |
      | prototype | git rebase main             |
      |           | git rebase origin/prototype |
    And the current branch is still "prototype"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE            |
      | main      | local, origin | main local commit  |
      |           |               | main origin commit |
      | prototype | local, origin | origin commit      |
      |           | local         | main local commit  |
      |           |               | main origin commit |
      |           |               | local commit       |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                                              |
      | prototype | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "prototype"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE            |
      | main      | local, origin | main local commit  |
      |           |               | main origin commit |
      | prototype | local         | local commit       |
      |           | origin        | origin commit      |
    And the initial branches and lineage exist
