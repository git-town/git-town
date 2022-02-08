Feature: rename a perennial branch

  Background:
    Given my repo has a perennial branch "production"
    And my repo has a feature branch "hotfix" as a child of "production"
    And the commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | hotfix     | local, remote | hotfix commit     |
      | production | local, remote | production commit |
    And I am on the "production" branch

  Scenario: normal rename fails
    When I run "git-town rename-branch production new"
    Then it runs no commands
    And it prints the error:
      """
      "production" is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'
      """

  Scenario: forced rename works
    When I run "git-town rename-branch --force production new"
    Then it runs the commands
      | BRANCH     | COMMAND                     |
      | production | git fetch --prune --tags    |
      |            | git branch new production   |
      |            | git checkout new            |
      | new        | git push -u origin new      |
      |            | git push origin :production |
      |            | git branch -D production    |
    And I am now on the "new" branch
    And the perennial branches are now "new"
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE           |
      | hotfix | local, remote | hotfix commit     |
      | new    | local, remote | production commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | hotfix | new    |

  Scenario: undo
    Given I ran "git-town rename-branch --force production new"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                                             |
      | new        | git branch production {{ sha 'production commit' }} |
      |            | git push -u origin production                       |
      |            | git push origin :new                                |
      |            | git checkout production                             |
      | production | git branch -D new                                   |
    And I am now on the "production" branch
    And the perennial branches are now "production"
    And now the initial commits exist
    And my repo now has its initial branches and branch hierarchy
