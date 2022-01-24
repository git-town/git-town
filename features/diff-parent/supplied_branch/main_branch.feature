Feature: git town-diff-parent: errors when trying to diff the main branch

  (see ../current_branch/on_main_branch.feature)

  Scenario: result
    Given my repo has a feature branch named "feature"
    And I am on the "feature" branch
    When I run "git-town diff-parent main"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """
    And I am still on the "feature" branch
