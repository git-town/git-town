@skipWindows
Feature: prompt for parent branch when unknown

  Scenario Outline:
    Given my repo has a branch "feature-1"
    And I am on the "feature-1" branch
    When I run "git-town <COMMAND>" and answer the prompts:
      | PROMPT                                          | ANSWER  |
      | Please specify the parent branch of 'feature-1' | [ENTER] |

    Examples:
      | COMMAND           |
      | append feature-2  |
      | diff-parent       |
      | kill feature-1    |
      | prepend feature-2 |
      | sync              |

  Scenario: prompt for parent branch when running git town-sync --all
    Given my repo has a branch "feature-1"
    And my repo has a branch "feature-2"
    And my repo contains the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, remote | main commit      |
      | feature-1 | local, remote | feature-1 commit |
      | feature-2 | local, remote | feature-2 commit |
    And I am on the "main" branch
    When I run "git-town sync --all" and answer the prompts:
      | PROMPT                                          | ANSWER  |
      | Please specify the parent branch of 'feature-1' | [ENTER] |
      | Please specify the parent branch of 'feature-2' | [ENTER] |
    Then my repo now has the commits
      | BRANCH    | LOCATION      | MESSAGE                            |
      | main      | local, remote | main commit                        |
      | feature-1 | local, remote | feature-1 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-1 |
      | feature-2 | local, remote | feature-2 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-2 |
