Feature: git town-hack: starting a new feature from a feature branch (with remote repo)

  As a developer working on something unrelated to my current feature branch
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.


  Background:
    Given I have a feature branch named "existing-feature"
    And the following commits exist in my repository
      | BRANCH           | LOCATION | MESSAGE                 |
      | main             | remote   | main commit             |
      | existing-feature | local    | existing feature commit |
    And I am on the "existing-feature" branch
    And I have an uncommitted file
    When I run `gt hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                          |
      | existing-feature | git fetch --prune                |
      |                  | git add -A                       |
      |                  | git stash                        |
      |                  | git checkout main                |
      | main             | git rebase origin/main           |
      |                  | git checkout -b new-feature main |
      | new-feature      | git push -u origin new-feature   |
      |                  | git stash pop                    |
    And I end up on the "new-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH           | LOCATION         | MESSAGE                 |
      | main             | local and remote | main commit             |
      | existing-feature | local            | existing feature commit |
      | new-feature      | local and remote | main commit             |
