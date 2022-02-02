Feature: shipping the supplied feature branch from a subfolder

  Background:
    Given my repo has the feature branches "feature" and "other-feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | remote   | feature commit |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file with name "new_folder/other_feature_file" and content "other feature content"
    When I run "git-town ship feature -m 'feature done'" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune --tags           |
      |               | git add -A                         |
      |               | git stash                          |
      |               | git checkout main                  |
      | main          | git rebase origin/main             |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      |               | git merge --no-edit main           |
      |               | git checkout main                  |
      | main          | git merge --squash feature         |
      |               | git commit -m "feature done"       |
      |               | git push                           |
      |               | git push origin :feature           |
      |               | git branch -D feature              |
      |               | git checkout other-feature         |
      | other-feature | git stash pop                      |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, remote | feature done |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | other-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH        | COMMAND                                       |
      | other-feature | git add -A                                    |
      |               | git stash                                     |
      |               | git checkout main                             |
      | main          | git branch feature {{ sha 'feature commit' }} |
      |               | git push -u origin feature                    |
      |               | git revert {{ sha 'feature done' }}           |
      |               | git push                                      |
      |               | git checkout feature                          |
      | feature       | git reset --hard {{ sha 'Initial commit' }}   |
      |               | git checkout main                             |
      | main          | git checkout other-feature                    |
      | other-feature | git stash pop                                 |
    And I am now on the "other-feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, remote | feature done          |
      |         |               | Revert "feature done" |
      | feature | remote        | feature commit        |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | feature       | main   |
      | other-feature | main   |
