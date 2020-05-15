Feature: git town-diff-parent: errors when trying to diff the main branch

  As a developer accidentally trying to diff the main branch
  I should see an error that I cannot diff the main branch
  Because the master branch cannot have a parent branch


  Scenario: result
    Given my repository has a feature branch named "feature"
    And I am on the "main" branch
    When I run "git-town diff-parent"
    Then it runs no commands
    And it prints the error:
      """
      You can only diff-parent feature branches
      """
    And I am still on the "main" branch
