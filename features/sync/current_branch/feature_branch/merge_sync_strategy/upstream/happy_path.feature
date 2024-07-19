Feature: with upstream repo

  Background:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And an upstream repo
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | main    | upstream | upstream commit |
      | feature | local    | local commit    |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git fetch upstream main                 |
      |         | git rebase upstream/main                |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
      |         | git merge --no-edit --ff main           |
      |         | git push                                |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE                          |
      | main    | local, origin, upstream | upstream commit                  |
      | feature | local, origin           | local commit                     |
      |         |                         | upstream commit                  |
      |         |                         | Merge branch 'main' into feature |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                               |
      | feature | git reset --hard {{ sha 'local commit' }}                             |
      |         | git push --force-with-lease origin {{ sha 'initial commit' }}:feature |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION                | MESSAGE         |
      | main    | local, origin, upstream | upstream commit |
      | feature | local                   | local commit    |
    And the initial branches and lineage exist
