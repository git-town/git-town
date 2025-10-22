@messyoutput @skipWindows
Feature: ship a coworker's feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE          | AUTHOR                            |
      | feature | local, origin | developer commit | developer <developer@example.com> |
      |         |               | coworker commit  | coworker <coworker@example.com>   |
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "feature"
    And I ran "git config --global --unset user.email"
    And I ran "git config --global --unset user.name"

  @this
  Scenario: choose the account configured by the GIT_AUTHOR_NAME and GIT_AUTHOR_EMAIL env variables
    When I run "git-town ship -m 'feature done'" with the environment variables "GIT_AUTHOR_NAME=developer" and "GIT_AUTHOR_EMAIL=developer@example.com" and enter into the dialog:
      | DIALOG               | KEYS       |
      | squash commit author | down enter |
    Then Git Town prints:
      """
      xxx
      """
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m "feature done"    |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
    And no lineage exists now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                            |
      | main   | local, origin | feature done | developer <developer@example.com> |

  Scenario: choose the account configured by the GIT_COMMITTER_NAME and GIT_COMMITTER_EMAIL env variables
    When I run "git-town ship -m 'feature done'" with the environment variables "GIT_COMMITTER_NAME=developer" and "GIT_COMMITTER_EMAIL=developer@example.com" and enter into the dialog:
      | DIALOG               | KEYS       |
      | squash commit author | down enter |
    Then Git Town runs the commands
      | BRANCH  | COMMAND                         |
      | feature | git fetch --prune --tags        |
      |         | git checkout main               |
      | main    | git merge --squash --ff feature |
      |         | git commit -m "feature done"    |
      |         | git push                        |
      |         | git push origin :feature        |
      |         | git branch -D feature           |
    And no lineage exists now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | feature done |
# NOTE: cannot verify the AUTHOR field of the commit here because Git uses the system user in this situation

  Scenario: no Git user configured
    When I run "git-town ship -m 'feature done'" and enter into the dialog:
      | DIALOG               | KEYS       |
      | squash commit author | down enter |
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                   |
      | feature | git fetch --prune --tags                                                  |
      |         | git checkout main                                                         |
      | main    | git merge --squash --ff feature                                           |
      |         | git commit -m "feature done" --author "developer <developer@example.com>" |
      |         | git push                                                                  |
      |         | git push origin :feature                                                  |
      |         | git branch -D feature                                                     |
    And no lineage exists now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                            |
      | main   | local, origin | feature done | developer <developer@example.com> |
