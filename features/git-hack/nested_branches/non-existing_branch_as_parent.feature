Feature: Trying to create a feature branch with a non-existing parent

  As a developer waiting for permission to ship a feature branch that contains changes needed for the next feature
  I want to be able to start working on the next feature while having access to the changes currently under review
  So that I am not slowed down by reviews and can keep working on my backlog as planned.


  Scenario: Creating a child branch of the current feature branch
    Given I have a feature branch named "feature"
    And Git Town knows that "feature" has the parent "main" and the parents "main"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    |
      | main    | local and remote | main_commit    | main_file    |
      | feature | local            | feature_commit | feature_file |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git hack child-feature zonk`
    Then I get the error
      """
      A branch named 'zonk' does not exist
      """
    And it runs no Git commands
    And I end up on the "feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
