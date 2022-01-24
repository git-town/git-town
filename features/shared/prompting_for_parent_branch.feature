@skipWindows
Feature: Prompt for parent branch when unknown


  Scenario: prompting for parent branch when running git town-append
    Given my repo has a branch "feature-1"
    And I am on the "feature-1" branch
    When I run "git-town append feature-2" and answer the prompts:
      | PROMPT                                          | ANSWER  |
      | Please specify the parent branch of 'feature-1' | [ENTER] |
    Then I am now on the "feature-2" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |


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


  Scenario: prompting for parent branch when running git town-kill
    Given my repo has a branch "feature"
    And I am on the "feature" branch
    When I run "git-town kill" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then I am now on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |


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


  Scenario: prompting for parent branch when running git town-sync
    Given my repo has a branch "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE        |
      | main    | local, remote | main commit    |
      | feature | local, remote | feature commit |
    And I am on the "feature" branch
    When I run "git-town sync" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, remote | main commit                      |
      | feature | local, remote | feature commit                   |
      |         |               | main commit                      |
      |         |               | Merge branch 'main' into feature |


  Scenario: prompting for parent branch when running git town-sync --all
    Given my repo has a branch "feature-1"
    And my repo has a branch "feature-2"
    And the following commits exist in my repo
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
