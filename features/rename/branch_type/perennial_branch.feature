Feature: rename a perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE      | PARENT     | LOCATIONS     |
      | production | perennial |            | local, origin |
      | hotfix     | feature   | production | local, origin |
    And the commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | hotfix     | local, origin | hotfix commit     |
      | production | local, origin | production commit |
    And the current branch is "production"

  Scenario: normal rename fails
    When I run "git-town rename production new"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And Git Town prints the error:
      """
      "production" is a perennial branch. Renaming a perennial branch typically requires other updates. If you are sure you want to do this, use '--force'
      """

  Scenario: forced rename works
    When I run "git-town rename --force production new"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                          |
      | production | git fetch --prune --tags         |
      |            | git branch --move production new |
      |            | git checkout new                 |
      | new        | git push -u origin new           |
      |            | git push origin :production      |
    And this lineage exists now
      """
      new
        hotfix
      """
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE           |
      | new    | local, origin | production commit |
      | hotfix | local, origin | hotfix commit     |
    And the perennial branches are now "new"

  Scenario: undo
    Given I run "git-town rename --force production new"
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                             |
      | new        | git branch production {{ sha 'production commit' }} |
      |            | git push -u origin production                       |
      |            | git checkout production                             |
      | production | git branch -D new                                   |
      |            | git push origin :new                                |
    And the initial branches and lineage exist now
    And the perennial branches are now "production"
    And the initial commits exist now
