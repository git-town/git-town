Feature: git town-diff-parent: errors when trying to diff the main branch

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given my repository has a feature branch named "feature"
    And I am on the "feature" branch


  Scenario: result
    When I run "git-town diff-parent main"
    Then it runs no commands
    And it prints the error:
      """
      You can only diff-parent feature branches
      """
    And I am still on the "feature" branch
