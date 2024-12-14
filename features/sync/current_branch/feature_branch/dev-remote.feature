Feature: sync a branch to a custom dev-remote

  Background: I fetch updates while resolving merge conflicts
    Given a Git repo with origin
    And I rename the "origin" remote to "fork"
    And Git Town setting "dev-remote" is "fork"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE |
      | branch | local    | commit  |
    And the current branch is "branch"
    And I run "git-town sync"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                               |
      | branch | git fetch --prune --tags              |
      |        | git checkout main                     |
      | main   | git rebase fork/main --no-update-refs |
      |        | git checkout branch                   |
      | branch | git merge --no-edit --ff main         |
      |        | git push -u fork branch               |
    And all branches are now synchronized
    And the current branch is still "branch"
    And these branches exist now
      | REPOSITORY    | BRANCHES     |
      | local, origin | main, branch |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE |
      | branch | local, origin | commit  |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND               |
      | branch | git push fork :branch |
    And the current branch is still "branch"
    And these branches exist now
      | REPOSITORY | BRANCHES     |
      | local      | main, branch |
      | origin     | main         |
    And these commits exist now
      | BRANCH | LOCATION | MESSAGE |
      | branch | local    | commit  |
