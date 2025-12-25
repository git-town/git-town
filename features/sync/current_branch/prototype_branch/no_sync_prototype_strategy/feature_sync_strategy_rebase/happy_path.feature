Feature: sync the current prototype branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS |
      | prototype | prototype | main   | local     |
    And the commits
      | BRANCH    | LOCATION | MESSAGE            |
      | main      | local    | main local commit  |
      |           | local    | main origin commit |
      | prototype | local    | local commit       |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "prototype"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                      |
      | prototype | git fetch --prune --tags                                                     |
      |           | git checkout main                                                            |
      | main      | git -c rebase.updateRefs=false rebase origin/main                            |
      |           | git push                                                                     |
      |           | git checkout prototype                                                       |
      | prototype | git -c rebase.updateRefs=false rebase --onto main {{ sha 'initial commit' }} |
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE            |
      | main      | local, origin | main local commit  |
      |           |               | main origin commit |
      | prototype | local         | local commit       |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                           |
      | prototype | git reset --hard {{ sha-initial 'local commit' }} |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE            |
      | main      | local, origin | main local commit  |
      |           |               | main origin commit |
      | prototype | local         | local commit       |
