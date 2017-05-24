Feature: push branch to remote upon creation

  Background:
    Given my repository has the "hack-push-flag" configuration set to "true"
    And I have a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE                 | FILE NAME             | FILE CONTENT             |
      | existing-feature | local and remote | existing_feature_commit | existing_feature_file | existing feature content |
    And I am on the "existing-feature" branch
    And I have an uncommitted file


  Scenario: inserting a branch into the branch ancestry
    When I run `git-town prepend new-parent`
    Then it runs the commands
      | BRANCH           | COMMAND                    |
      | existing-feature | git fetch --prune          |
      |                  | git add -A                 |
      |                  | git stash                  |
      |                  | git checkout main          |
      | main             | git rebase origin/main     |
      |                  | git branch new-parent main |
      |                  | git checkout new-parent    |
      | new-parent       | git push -u origin new-parent |
      | | git stash pop              |
    And I end up on the "new-parent" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE                 |
      | existing-feature | local and remote | existing_feature_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT     |
      | existing-feature | new-parent |
      | new-parent       | main       |


  Scenario: Undo
    Given I run `git-town prepend new-parent`
    When I run `git-town prepend --undo`
    Then it runs the commands
        | BRANCH           | COMMAND                       |
        | new-parent       | git add -A                    |
        |                  | git stash                     |
        |                  | git push origin :new-parent   |
        | | git checkout main             |
        | main             | git branch -d new-parent      |
        |                  | git checkout existing-feature |
        | existing-feature | git stash pop                 |
    And I end up on the "existing-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | existing-feature | main   |
