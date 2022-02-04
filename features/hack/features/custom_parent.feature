@skipWindows
Feature: customize the parent for the new feature branch

  Background:
    Given my repo has a branch "feature-1"
    And I am on the "feature-1" branch
    When I run "git-town hack -p feature-2" and answer the prompts:
      | PROMPT                                          | ANSWER        |
      | Please specify the parent branch of 'feature-2' | [DOWN][ENTER] |
      | Please specify the parent branch of 'feature-1' | [ENTER]       |

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                        |
      | feature-1 | git fetch --prune --tags       |
      |           | git merge --no-edit main       |
      |           | git push -u origin feature-1   |
      |           | git branch feature-2 feature-1 |
      |           | git checkout feature-2         |
    And I am now on the "feature-2" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                    |
      | feature-2 | git checkout feature-1     |
      | feature-1 | git branch -d feature-2    |
      |           | git push origin :feature-1 |
    And I am now on the "feature-1" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT |
      | feature-1 | main   |
