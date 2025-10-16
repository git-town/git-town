Feature: sync the current branch that is contribution via regex

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE   | LOCATIONS     |
      | renovate/1 | (none) | local, origin |
    And the commits
      | BRANCH     | LOCATION      | MESSAGE       | FILE NAME   |
      | main       | local, origin | main commit   | main_file   |
      | renovate/1 | local         | local commit  | local_file  |
      |            | origin        | origin commit | origin_file |
    And Git setting "git-town.contribution-regex" is "^renovate"
    And the current branch is "renovate/1"
    When I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                 |
      | renovate/1 | git fetch --prune --tags                                |
      |            | git -c rebase.updateRefs=false rebase origin/renovate/1 |
      |            | git push                                                |
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE       |
      | main       | local, origin | main commit   |
      | renovate/1 | local, origin | origin commit |
      |            |               | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                                           |
      | renovate/1 | git reset --hard {{ sha-initial 'local commit' }}                                 |
      |            | git push --force-with-lease origin {{ sha-in-origin 'origin commit' }}:renovate/1 |
    And the initial branches and lineage exist now
    And the initial commits exist now
