Feature: git kill: killing the given feature branch when on it (without open changes or remote origin)

  (see ./with_open_changes.feature)


  Background:
    Given I have feature branches named "feature" and "dead-feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH       | LOCATION | MESSAGE         | FILE NAME        |
      | feature      | local    | good commit     | good_file        |
      | dead-feature | local    | dead-end commit | unfortunate_file |
    And I am on the "dead-feature" branch
    When I run `git kill dead-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH       | COMMAND                    |
      | dead-feature | git checkout main          |
      | main         | git branch -D dead-feature |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
    And I have the following commits
      | BRANCH  | LOCATION | MESSAGE     | FILE NAME |
      | feature | local    | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                       |
      | main   | git branch dead-feature [SHA:dead-end commit] |
      | main   | git checkout dead-feature                     |
    And I end up on the "dead-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                    |
      | local      | main, dead-feature, feature |
    And I have the following commits
      | BRANCH       | LOCATION | MESSAGE         | FILE NAME        |
      | feature      | local    | good commit     | good_file        |
      | dead-feature | local    | dead-end commit | unfortunate_file |

