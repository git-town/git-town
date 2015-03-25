Feature: git kill: killing the given feature branch (without open changes)

  (see ./with_open_changes.feature)


  Background:
    Given I have feature branches named "current-feature" and "dead-feature"
    And the following commits exist in my repository
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME        |
      | current-feature | local and remote | current feature commit | good_file        |
      | dead-feature    | local and remote | dead-end commit        | unfortunate_file |
    And I am on the "current-feature" branch
    When I run `git kill dead-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH          | COMMAND                       |
      | current-feature | git fetch --prune             |
      | current-feature | git push origin :dead-feature |
      | current-feature | git branch -D dead-feature    |
    And I am still on the "current-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES              |
      | local      | main, current-feature |
      | remote     | main, current-feature |
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME |
      | current-feature | local and remote | current feature commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH          | COMMAND                                              |
      | current-feature | git branch dead-feature <%= sha 'dead-end commit' %> |
      | current-feature | git push -u origin dead-feature                      |
    And I am still on the "current-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                            |
      | local      | main, dead-feature, current-feature |
      | remote     | main, dead-feature, current-feature |
    And I have the following commits
      | BRANCH          | LOCATION         | MESSAGE                | FILE NAME        |
      | current-feature | local and remote | current feature commit | good_file        |
      | dead-feature    | local and remote | dead-end commit        | unfortunate_file |
