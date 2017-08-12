Feature: git append: offline mode

  When having no internet connection
  I want that new branches are created without attempting network accesses
  So that I don't see unnecessary errors.


  Background:
    Given Git Town is in offline mode
    And I have a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION         | MESSAGE                 |
      | existing-feature | local and remote | existing feature commit |
    And I am on the "existing-feature" branch
    And I have an uncommitted file


  Scenario: appending a branch in offline mode
    When I run `git-town append new-feature`
    Then it runs the commands
      | BRANCH           | COMMAND                                     |
      | existing-feature | git add -A                                  |
      |                  | git stash                                   |
      |                  | git checkout main                           |
      | main             | git rebase origin/main                      |
      |                  | git checkout existing-feature               |
      | existing-feature | git merge --no-edit origin/existing-feature |
      |                  | git merge --no-edit main                    |
      |                  | git branch new-feature existing-feature     |
      |                  | git checkout new-feature                    |
      | new-feature      | git stash pop                               |
    And I end up on the "new-feature" branch
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE                 |
      | existing-feature | local and remote | existing feature commit |
      | new-feature      | local            | existing feature commit |


  Scenario: Undo
    Given I run `git-town append new-feature`
    When I run `git-town append --undo`
    Then it runs the commands
        | BRANCH           | COMMAND                       |
        | new-feature      | git add -A                    |
        |                  | git stash                     |
        |                  | git checkout existing-feature |
        | existing-feature | git branch -D new-feature     |
        |                  | git checkout main             |
        | main             | git checkout existing-feature |
        | existing-feature | git stash pop                 |
    And I end up on the "existing-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is now aware of this branch hierarchy
      | BRANCH           | PARENT |
      | existing-feature | main   |
