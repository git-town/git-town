Feature: rename with configured branch-prefix via environment variable

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT | LOCATIONS     |
      | feature-1 | feature | main   | local, origin |
    And the current branch is "feature-1"

  Scenario Outline:
    When I run "git-town rename <BRANCH_NAME>" with these environment variables
      | GIT_TOWN_BRANCH_PREFIX | kg- |
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                  |
      | feature-1    | git fetch --prune --tags                 |
      |              | git branch --move feature-1 kg-feature-2 |
      |              | git checkout kg-feature-2                |
      | kg-feature-2 | git push -u origin kg-feature-2          |
      |              | git push origin :feature-1               |
    And the current branch is now "kg-feature-2"
    And this lineage exists now
      """
      main
        kg-feature-2
      """

    Examples:
      | BRANCH_NAME  |
      | feature-2    |
      | kg-feature-2 |

  Scenario: undo
    When I run "git-town rename feature-2" with these environment variables
      | GIT_TOWN_BRANCH_PREFIX | kg- |
    And I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH       | COMMAND                                         |
      | kg-feature-2 | git branch feature-1 {{ sha 'initial commit' }} |
      |              | git push -u origin feature-1                    |
      |              | git checkout feature-1                          |
      | feature-1    | git branch -D kg-feature-2                      |
      |              | git push origin :kg-feature-2                   |
    And the current branch is now "feature-1"
    And this lineage exists now
      """
      main
        feature-1
      """
