Feature: git town-ship: shipping a coworker's feature branch

  As a developer shipping a coworker's feature branch
  I want my coworker to be the author of the commit added to the main branch
  So that my coworker is given credit for their work


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE         | AUTHOR                            |
      | feature | local    | feature commit1 | developer <developer@example.com> |
      |         |          | feature commit2 | developer <developer@example.com> |
      |         |          | feature commit3 | coworker <coworker@example.com>   |
    And I am on the "feature" branch


  Scenario Outline: prompt for squashed commit author
    When I run "git-town ship -m 'feature done'" and answer the prompts:
      | PROMPT                                        | ANSWER   |
      | Please choose an author for the squash commit | <ANSWER> |
    And my repository now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR           |
      | main   | local, remote | feature done | <FEATURE_AUTHOR> |

    Examples:
      | ANSWER        | FEATURE_AUTHOR                    |
      | [ENTER]       | developer <developer@example.com> |
      | [DOWN][ENTER] | coworker <coworker@example.com>   |
