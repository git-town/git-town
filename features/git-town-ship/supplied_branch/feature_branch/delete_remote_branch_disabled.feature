Feature: Skip deleting the remote branch when shipping another branch

  When using GitHub's feature to automatically delete head branches of pull requests.
  I want "git ship" to skip deleting the remote feature branch
  So that I can keep using Git Town in this situation.


  Background:
    Given my repository has the feature branches "feature" and "other-feature"
    And the following commits exist in my repository
      | BRANCH        | LOCATION      | MESSAGE        | FILE NAME    |
      | feature       | local, remote | feature commit | feature_file |
      | other-feature | local         | other commit   | other_file   |
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
    And I end up on the "other-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And my repository now has the following commits
      | BRANCH        | LOCATION      | MESSAGE      | FILE NAME    |
      | main          | local, remote | feature done | feature_file |
      | other-feature | local         | other commit | other_file   |


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
    And I end up on the "other-feature" branch
    And my repository now has the following commits
      | BRANCH        | LOCATION      | MESSAGE               | FILE NAME    |
      | main          | local, remote | feature done          | feature_file |
      |               |               | Revert "feature done" | feature_file |
      | feature       | local         | feature commit        | feature_file |
      | other-feature | local         | other commit          | other_file   |
