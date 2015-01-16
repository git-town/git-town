Feature: git kill: current feature branch with a tracking branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have feature branches named "feature" and "dead-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | feature      | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
    And I am on the "dead-feature" branch
    When I run `git kill`


  Scenario: result
    Then it runs the Git commands
      | BRANCH       | COMMAND                       |
      | dead-feature | git fetch --prune             |
      | dead-feature | git checkout main             |
      | main         | git push origin :dead-feature |
      | main         | git branch -D dead-feature    |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE     | FILE NAME |
      | feature | local and remote | good commit | good_file |


  Scenario: Undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                       |
      | main   | git branch dead-feature [SHA:dead-end commit] |
      | main   | git push -u origin dead-feature               |
      | main   | git checkout dead-feature                     |
    And I end up on the "dead-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                    |
      | local      | main, dead-feature, feature |
      | remote     | main, dead-feature, feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | feature      | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
