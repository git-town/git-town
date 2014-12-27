Feature: Git Kill: killing the current feature by name, without open changes

  Background:
    Given I have feature branches named "good-feature" and "dead-feature"
    And the following commits exist in my repository
      | BRANCH       | LOCATION         | MESSAGE         | FILE NAME        |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |
    And I am on the "dead-feature" branch
    When I run `git kill dead-feature`


  Scenario: result
    Then it runs the Git commands
      | BRANCH       | COMMAND                       |
      | dead-feature | git fetch --prune             |
      | dead-feature | git checkout main             |
      | main         | git push origin :dead-feature |
      | main         | git branch -D dead-feature    |
    And I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES           |
      | local      | main, good-feature |
      | remote     | main, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE     | FILES     |
      | good-feature | local and remote | good commit | good_file |


  Scenario: undoing the kill
    When I run `git kill --undo`
    Then it runs the Git commands
      | BRANCH | COMMAND                                       |
      | main   | git branch dead-feature [SHA:dead-end commit] |
      | main   | git push -u origin dead-feature               |
      | main   | git checkout dead-feature                     |
    And I end up on the "dead-feature" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                         |
      | local      | main, dead-feature, good-feature |
      | remote     | main, dead-feature, good-feature |
    And I have the following commits
      | BRANCH       | LOCATION         | MESSAGE         | FILES            |
      | good-feature | local and remote | good commit     | good_file        |
      | dead-feature | local and remote | dead-end commit | unfortunate_file |

