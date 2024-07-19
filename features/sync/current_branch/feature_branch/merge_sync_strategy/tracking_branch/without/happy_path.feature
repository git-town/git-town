Feature: sync the current feature branch without a tracking branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      |         | origin   | origin main commit   |
      | feature | local    | local feature commit |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git fetch --prune --tags      |
      |         | git checkout main             |
      | main    | git rebase origin/main        |
      |         | git push                      |
      |         | git checkout feature          |
      | feature | git merge --no-edit --ff main |
      |         | git push -u origin feature    |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | origin main commit               |
      |         |               | local main commit                |
      | feature | local, origin | local feature commit             |
      |         |               | origin main commit               |
      |         |               | local main commit                |
      |         |               | Merge branch 'main' into feature |
    And the branches are now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                           |
      | feature | git push origin :feature                          |
      |         | git reset --hard {{ sha 'local feature commit' }} |
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE              |
      | main    | local, origin | origin main commit   |
      |         |               | local main commit    |
      | feature | local         | local feature commit |
    And the initial branches and lineage exist
