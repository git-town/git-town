Feature: push branch to remote upon creation

  Background:
    Given the new-branch-push-flag configuration is true
    And my repo has a feature branch named "existing-feature"
    And the following commits exist in my repo
      | BRANCH           | LOCATION      | MESSAGE                 | FILE NAME             | FILE CONTENT             |
      | existing-feature | local, remote | existing_feature_commit | existing_feature_file | existing feature content |
    And I am on the "existing-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town prepend new-parent"

  Scenario: inserting a branch into the branch ancestry
    Then it runs the commands
      | BRANCH           | COMMAND                       |
      | existing-feature | git fetch --prune --tags      |
      |                  | git add -A                    |
      |                  | git stash                     |
      |                  | git checkout main             |
      | main             | git rebase origin/main        |
      |                  | git branch new-parent main    |
      |                  | git checkout new-parent       |
      | new-parent       | git push -u origin new-parent |
      |                  | git stash pop                 |
    And I am now on the "new-parent" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH           | LOCATION      | MESSAGE                 |
      | existing-feature | local, remote | existing_feature_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT     |
      | existing-feature | new-parent |
      | new-parent       | main       |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH           | COMMAND                       |
      | new-parent       | git add -A                    |
      |                  | git stash                     |
      |                  | git push origin :new-parent   |
      |                  | git checkout main             |
      | main             | git branch -d new-parent      |
      |                  | git checkout existing-feature |
      | existing-feature | git stash pop                 |
    And I am now on the "existing-feature" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | existing-feature | main   |
