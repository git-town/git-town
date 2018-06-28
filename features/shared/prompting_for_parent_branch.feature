Feature: Prompt for parent branch when unknown

  As a developer running a command on a branch without a parent branch
  I should see a prompt asking for the information
  So the command can work as I expect


  Scenario: prompting for parent branch when running git town-append
    Given my repository has a feature branch named "feature-1" with no parent
    And I am on the "feature-1" branch
    When I run `git-town append feature-2` and answer the prompts:
      | PROMPT                                          | ANSWER  |
      | Please specify the parent branch of 'feature-1' | [ENTER] |
    Then I end up on the "feature-2" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |


  Scenario: prompting for parent branch when running git town-hack -p
    Given my repository has a feature branch named "feature-1" with no parent
    And I am on the "feature-1" branch
    When I run `git-town hack -p feature-2` and answer the prompts:
      | PROMPT                                          | ANSWER        |
      | Please specify the parent branch of 'feature-2' | [DOWN][ENTER] |
      | Please specify the parent branch of 'feature-1' | [ENTER]       |
    Then I end up on the "feature-2" branch
    And Git Town is now aware of this branch hierarchy
      | BRANCH    | PARENT    |
      | feature-1 | main      |
      | feature-2 | feature-1 |


  Scenario: prompting for parent branch when running git town-kill
    Given my repository has a feature branch named "feature" with no parent
    And I am on the "feature" branch
    When I run `git-town kill` and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |


  Scenario: prompting for parent branch when running git town-new-pull-request
    Given I have "open" installed
    And my repository has a feature branch named "feature"
    And Git Town has no branch hierarchy information for "feature"
    And my repo's remote origin is git@github.com:Originate/git-town.git
    And I am on the "feature" branch
    When I run `git-town new-pull-request` and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then I see a new GitHub pull request for the "feature" branch in the "Originate/git-town" repo in my browser


  Scenario: prompting for parent branch when running git town-sync
    Given my repository has a feature branch named "feature" with no parent
    And the following commits exist in my repository
      | BRANCH  | LOCATION         | MESSAGE        |
      | main    | local and remote | main commit    |
      | feature | local and remote | feature commit |
    And I am on the "feature" branch
    When I run `git-town sync` and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then my repository has the following commits
      | BRANCH  | LOCATION         | MESSAGE                          |
      | main    | local and remote | main commit                      |
      | feature | local and remote | feature commit                   |
      |         |                  | main commit                      |
      |         |                  | Merge branch 'main' into feature |


  Scenario: prompting for parent branch when running git town-sync --all
    Given my repository has a feature branch named "feature-1" with no parent
    And my repository has a feature branch named "feature-2" with no parent
    And the following commits exist in my repository
      | BRANCH    | LOCATION         | MESSAGE          |
      | main      | local and remote | main commit      |
      | feature-1 | local and remote | feature-1 commit |
      | feature-2 | local and remote | feature-2 commit |
    And I am on the "main" branch
    When I run `git-town sync --all` and answer the prompts:
      | PROMPT                                          | ANSWER  |
      | Please specify the parent branch of 'feature-1' | [ENTER] |
      | Please specify the parent branch of 'feature-2' | [ENTER] |
    Then my repository has the following commits
      | BRANCH    | LOCATION         | MESSAGE                            |
      | main      | local and remote | main commit                        |
      | feature-1 | local and remote | feature-1 commit                   |
      |           |                  | main commit                        |
      |           |                  | Merge branch 'main' into feature-1 |
      | feature-2 | local and remote | feature-2 commit                   |
      |           |                  | main commit                        |
      |           |                  | Merge branch 'main' into feature-2 |
