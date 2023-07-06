Feature: ship the supplied feature branch from a subfolder

  Background:
    Given the feature branches "feature" and "other"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | origin   | feature commit |
    And the current branch is "other"
    And an uncommitted file with name "new_folder/other_feature_file" and content "other feature content"
    When I run "git-town ship feature -m 'feature done'" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | other   | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
      |         | git checkout other                 |
      | other   | git stash pop                      |
    And the current branch is now "other"
    And the uncommitted file still exists
    And the branches are now
      | REPOSITORY    | BRANCHES    |
      | local, origin | main, other |
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      |
      | main   | local, origin | feature done |
    And this branch lineage exists now
      | BRANCH | PARENT |
      | other  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | other   | git add -A                                    |
      |         | git stash                                     |
      |         | git checkout main                             |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git push -u origin feature                    |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git push                                      |
      |         | git checkout feature                          |
      | feature | git reset --hard {{ sha 'Initial commit' }}   |
      |         | git checkout main                             |
      | main    | git checkout other                            |
      | other   | git stash pop                                 |
    And the current branch is now "other"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | origin        | feature commit        |
    And the initial branches and hierarchy exist
