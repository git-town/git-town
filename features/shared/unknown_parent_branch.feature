@skipWindows
Feature: Prompt for parent branch when unknown

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

  Scenario: prompting for parent branch when running git town-hack -p
    Given my repo has a branch "feature-1"
    And I am on the "feature-1" branch
    When I run "git-town hack -p feature-2" and answer the prompts:
      | PROMPT                                          | ANSWER        |
      | Please specify the parent branch of 'feature-2' | [DOWN][ENTER] |
      | Please specify the parent branch of 'feature-1' | [ENTER]       |
    Then I am now on the "feature-2" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |

  @skipWindows
  Scenario: prompting for parent branch when running git town-new-pull-request
    And my computer has the "open" tool installed
    And my repo has a branch "feature"
    And my repo's origin is "git@github.com:git-town/git-town.git"
    And I am on the "feature" branch
    When I run "git-town new-pull-request" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/feature?expand=1
      """

  Scenario: prompting for parent branch when running git town-sync --all
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
    Then my repo now has the following commits
      | BRANCH    | LOCATION      | MESSAGE                            |
      | main      | local, remote | main commit                        |
      | feature-1 | local, remote | feature-1 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-1 |
      | feature-2 | local, remote | feature-2 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-2 |
