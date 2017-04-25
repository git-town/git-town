Feature: git town-kill: errors when trying to kill the main branch

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     |
      | main    | local and remote | main commit |
      | feature | local and remote | good commit |
    And I am on the "feature" branch


  Scenario: result
    Given I have an uncommitted file
    When I run `git-town kill main`
    Then it runs no commands
    And I get the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And I still have my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I am left with my original commits
