Feature: git prepend: offline mode

  When having no internet connection
  I want that new branches are created without attempting network accesses
  So that I don't see unnecessary errors.


  Background:
    Given Git Town is in offline mode
    And my repository has a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE                 | FILE NAME             | FILE CONTENT             |
      | existing-feature | local and remote | existing_feature_commit | existing_feature_file | existing feature content |
    And I am on the "existing-feature" branch
    And my workspace has an uncommitted file


  Scenario: inserting a branch into the branch ancestry
    When I run `git-town prepend new-parent`
    Then it runs the commands
      | BRANCH           | COMMAND                    |
      | existing-feature | git add -A                 |
      |                  | git stash                  |
      |                  | git checkout main          |
      | main             | git rebase origin/main     |
      |                  | git branch new-parent main |
      |                  | git checkout new-parent    |
      | new-parent       | git stash pop              |
    And I end up on the "new-parent" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH           | LOCATION         | MESSAGE                 |
      | existing-feature | local and remote | existing_feature_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT     |
      | existing-feature | new-parent |
      | new-parent       | main       |


  Scenario: Undo
    Given I run `git-town prepend new-parent`
    When I run `git-town undo`
    Then it runs the commands
        | BRANCH           | COMMAND                       |
        | new-parent       | git add -A                    |
        |                  | git stash                     |
        |                  | git checkout main             |
        | main             | git branch -d new-parent      |
        |                  | git checkout existing-feature |
        | existing-feature | git stash pop                 |
    And I end up on the "existing-feature" branch
    And my workspace still contains my uncommitted file
    And my repository is left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | existing-feature | main   |
