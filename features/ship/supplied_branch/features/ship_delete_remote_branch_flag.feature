Feature: skip deleting the remote branch when shipping another branch

  Background:
    Given my repo has the feature branches "feature" and "other"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
      | other   | local         | other commit   |
    And I am on the "other" branch
    And my repo has "git-town.ship-delete-remote-branch" set to "false"
    When I run "git-town ship feature -m 'feature done'"
    And the remote deletes the "feature" branch

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | other   | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git branch -D feature              |
      |         | git checkout other                 |
    And I am now on the "other" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES    |
      | local, remote | main, other |
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, remote | feature done |
      | other  | local         | other commit |
    And Git Town now knows about this branch hierarchy
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | other   | git checkout main                             |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git push                                      |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git checkout other                            |
    And I am now on the "other" branch
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, remote | feature done          |
      |         |               | Revert "feature done" |
      | feature | local         | feature commit        |
      | other   | local         | other commit          |
    And the existing branches are
      | REPOSITORY | BRANCHES             |
      | local      | main, feature, other |
      | remote     | main, other          |
    And Git Town now knows the initial branch hierarchy
