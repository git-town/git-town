Feature: rename a perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | PARENT     | LOCATIONS     |
      | production | perennial |            | local, origin |
      | hotfix     | feature   | production | local, origin |
    And the current branch is "production"
    And the commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | hotfix     | local, origin | hotfix commit     |
      | production | local, origin | production commit |

  Scenario: normal rename fails
    When I run "git-town rename production new"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      "production" is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'
      """

  Scenario: forced rename works
    When I run "git-town rename --force production new"
    Then it runs the commands
      | BRANCH     | COMMAND                          |
      | production | git fetch --prune --tags         |
      |            | git branch --move production new |
      |            | git checkout new                 |
      | new        | git push -u origin new           |
      |            | git push origin :production      |
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
    Given I ran "git-town rename --force production new"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH     | COMMAND                                             |
      | new        | git branch production {{ sha 'production commit' }} |
      |            | git push -u origin production                       |
      |            | git checkout production                             |
      | production | git branch -D new                                   |
      |            | git push origin :new                                |
    And the current branch is now "production"
    And the perennial branches are now "production"
    And the initial commits exist now
    And the initial branches and lineage exist now
