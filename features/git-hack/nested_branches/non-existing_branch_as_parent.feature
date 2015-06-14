Feature: Trying to create a feature branch with a non-existing parent

  As a developer trying to fork off a non-existing branch
  I want to get a reminder about my mistake
  So that I can try again with the correct parent branch name.


  Scenario: Creating a child branch off the current feature branch
    Given I have a feature branch named "feature"
    And Git Town knows that "feature" has the parent "main" and the parents "main"
    And the following commit exists in my repository
      | BRANCH | LOCATION         | MESSAGE     | FILE NAME |
      | main   | local and remote | main_commit | main_file |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git hack feature zonk`
    Then I get the error
      """
      A branch named 'zonk' does not exist
      """
    And it runs no Git commands
    And I end up on the "main" branch
    And I still have my uncommitted file
    And I am left with my original commits
    And Git Town is aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
