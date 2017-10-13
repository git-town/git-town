Feature: git town-ship: shipping the supplied feature branch from a subfolder

  (see ../../current_branch/on_feature_branch/without_open_changes/in_subfolder.feature)


  Background:
    Given my repository has the feature branches "feature" and "other-feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature commit | feature_file | feature content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file with name: "new_folder/other_feature_file" and content: "other feature content"
    When I run `git-town ship feature -m "feature done"` in the "new_folder" folder


  Scenario: result
    Then Git Town runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune                  |
      | <none>        | cd <%= git_root_folder %>          |
      | other-feature | git add -A                         |
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
      | <none>        | cd <%= git_folder "new_folder" %>  |
    And I end up on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no "feature" branch
    And my repository has the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
