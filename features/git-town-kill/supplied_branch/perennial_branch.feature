Feature: git town-kill: errors when trying to kill a perennial branch

  (see ../current_branch/on_perennial_branch.feature)


  Background:
    Given I have a feature branch named "feature"
    And I have a perennial branch named "qa"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     |
      | feature | local and remote | good commit |
      | qa      | local and remote | qa commit   |
    And I am on the "feature" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git-town kill qa`
    Then it runs no commands
    And I get the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES          |
      | local      | main, qa, feature |
      | remote     | main, qa, feature |
    And I am left with my original commits
