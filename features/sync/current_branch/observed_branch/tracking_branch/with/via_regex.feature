Feature: sync the current branch that is observed via regex

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE   | LOCATIONS     |
      | renovate/1 | (none) | local, origin |
    And the current branch is "renovate/1"
    And Git Town setting "observed-regex" is "^renovate"
    And the commits
      | BRANCH     | LOCATION      | MESSAGE       | FILE NAME   |
      | main       | local, origin | main commit   | main_file   |
      | renovate/1 | local         | local commit  | local_file  |
      |            | origin        | origin commit | origin_file |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                      |
      | renovate/1 | git fetch --prune --tags     |
      |            | git rebase origin/renovate/1 |
    And the current branch is still "renovate/1"
    And these commits exist now
      | BRANCH     | LOCATION      | MESSAGE       |
      | main       | local, origin | main commit   |
      | renovate/1 | local, origin | origin commit |
      |            | local         | local commit  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                                              |
      | renovate/1 | git reset --hard {{ sha-before-run 'local commit' }} |
    And the current branch is still "renovate/1"
    And the initial commits exist now
    And the initial branches and lineage exist now
