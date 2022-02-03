Feature: Skip deleting the remote branch when shipping another branch

  Background:
    Given my repo has the feature branches "feature" and "other-feature"
    And my repo contains the commits
      | BRANCH        | LOCATION      | MESSAGE        |
      | feature       | local, remote | feature commit |
      | other-feature | local         | other commit   |
    And I am on the "other-feature" branch
    And my repo has "git-town.ship-delete-remote-branch" set to "false"
    When I run "git-town ship feature -m 'feature done'"
    And the remote deletes the "feature" branch

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune --tags           |
      |               | git checkout main                  |
      | main          | git rebase origin/main             |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      |               | git merge --no-edit main           |
      |               | git checkout main                  |
      | main          | git merge --squash feature         |
      |               | git commit -m "feature done"       |
      |               | git push                           |
      |               | git branch -D feature              |
      |               | git checkout other-feature         |
    And I am now on the "other-feature" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES            |
      | local, remote | main, other-feature |
    And my repo now has the following commits
      | BRANCH        | LOCATION      | MESSAGE      |
      | main          | local, remote | feature done |
      | other-feature | local         | other commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | other-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH        | COMMAND                                       |
      | other-feature | git checkout main                             |
      | main          | git branch feature {{ sha 'feature commit' }} |
      |               | git revert {{ sha 'feature done' }}           |
      |               | git push                                      |
      |               | git checkout feature                          |
      | feature       | git checkout main                             |
      | main          | git checkout other-feature                    |
    And I am now on the "other-feature" branch
    And my repo now has the following commits
      | BRANCH        | LOCATION      | MESSAGE               |
      | main          | local, remote | feature done          |
      |               |               | Revert "feature done" |
      | feature       | local         | feature commit        |
      | other-feature | local         | other commit          |
    And Git Town now has the original branch hierarchy
