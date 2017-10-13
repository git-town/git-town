Feature: git town-kill: errors when trying to kill the main branch

  (see ../current_branch/on_main_branch.feature)


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE     |
      | main    | local and remote | main commit |
      | feature | local and remote | good commit |
    And I am on the "feature" branch


  Scenario: result
    Given my workspace has an uncommitted file
    When I run `git-town kill main`
    Then Git Town runs no commands
    And it prints the error "You can only kill feature branches"
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And my repository is left with my original commits
