Feature: git town-kill: errors when trying to kill a perennial branch

  (see ../current_branch/on_perennial_branch.feature)


  Background:
    Given my repository has a feature branch named "feature"
    And my repository has a perennial branch named "qa"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     |
      | feature | local and remote | good commit |
      | qa      | local and remote | qa commit   |
    And I am on the "feature" branch


  Scenario: result
    Given my workspace has an uncommitted file
    When I run `git-town kill qa`
    Then Git Town runs no commands
    And it prints the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES          |
      | local      | main, qa, feature |
      | remote     | main, qa, feature |
    And my repository is left with my original commits
