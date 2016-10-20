Feature: git town-ship: shipping a coworker's feature branch

  As a developer shipping a coworker's feature branch
  I want my coworker to be the author of the commit added to the main branch
  So that my coworker is given credit for their work


  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | AUTHOR                            |
      | feature | local    | feature commit1 | developer <developer@example.com> |
      |         |          | feature commit2 | developer <developer@example.com> |
      |         |          | feature commit3 | coworker <coworker@example.com>   |
    And I am on the "feature" branch


  Scenario Outline: prompt for squashed commit author
    When I run `git town-ship -m 'feature done'` and <ACTION>
    Then I see
      """
      Multiple people authored the 'feature' branch.
      Please choose an author for the squash commit.

        1: developer <developer@example.com> (2 commits)
        2: coworker <coworker@example.com> (1 commit)

      Enter user's number or a custom author (default: 1):
      """
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | AUTHOR           |
      | main   | local and remote | feature done | <FEATURE_AUTHOR> |

    Examples:
      | ACTION                             | FEATURE_AUTHOR                    |
      | press ENTER                        | developer <developer@example.com> |
      | enter "1"                          | developer <developer@example.com> |
      | enter "2"                          | coworker <coworker@example.com>   |
      | enter "other <other@example.com>"" | other <other@example.com>         |


  Scenario Outline: enter invalid number then valid number
    When I run `git town-ship -m 'feature done'` and enter "<NUMBER>" and "1"
    Then I see "error: invalid number"
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | AUTHOR                            |
      | main   | local and remote | feature done | developer <developer@example.com> |

    Examples:
      | NUMBER |
      | 0      |
      | 3      |


  Scenario: enter invalid custom author
    When I run `git town-ship -m 'feature done'` and enter "invalid"
    Then I get the error "Ship aborted because commit exited with error"
    And I am left with my original commits


  Scenario: supplying the author via command line arguments
    When I run `git town-ship -m 'feature done' --author='other <other@example.com>'`
    Then I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | AUTHOR                    |
      | main   | local and remote | feature done | other <other@example.com> |
