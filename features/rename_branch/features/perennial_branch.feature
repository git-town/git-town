Feature: rename a perennial branch

  Background:
    Given the current branch is a perennial branch "production"
    And a feature branch "hotfix" as a child of "production"
    And the commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | hotfix     | local, origin | hotfix commit     |
      | production | local, origin | production commit |

  Scenario: normal rename fails
    When I run "git-town rename-branch production new"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      "production" is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'
      """

  @debug @this
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
    And the current branch is now "new"
    And the perennial branches are now "new"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE           |
      | hotfix | local, origin | hotfix commit     |
      | new    | local, origin | production commit |
    And this lineage exists now
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
    And the current branch is now "production"
    And the perennial branches are now "production"
    And the initial commits exist
    And the initial branches and lineage exist
