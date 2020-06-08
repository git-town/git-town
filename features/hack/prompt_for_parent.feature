Feature: git town-hack: prompt for parent branch

  As a developer working on a new feature that should be based off another branch
  I want to be able to create a new up-to-date branch without checking out the parent
  So that I have a unified way for creating new branches


  Scenario: selecting the default branch (the main development branch)
    Given the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE     |
      | main   | remote   | main_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack -p new-feature" and answer the prompts:
      | PROMPT                                            | ANSWER  |
      | Please specify the parent branch of 'new-feature' | [ENTER] |
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git fetch --prune --tags    |
      |             | git add -A                  |
      |             | git stash                   |
      |             | git rebase origin/main      |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I end up on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE     |
      | main        | local, remote | main_commit |
      | new-feature | local         | main_commit |


  Scenario: selecting another branch
    Given my repo has the perennial branch "production"
    And the following commits exist in my repo
      | BRANCH     | LOCATION | MESSAGE           |
      | production | remote   | production_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack -p hotfix" and answer the prompts:
      | PROMPT                                       | ANSWER        |
      | Please specify the parent branch of 'hotfix' | [DOWN][ENTER] |
    Then it runs the commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune --tags     |
      |            | git add -A                   |
      |            | git stash                    |
      |            | git checkout production      |
      | production | git rebase origin/production |
      |            | git branch hotfix production |
      |            | git checkout hotfix          |
      | hotfix     | git stash pop                |
    And I end up on the "hotfix" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | hotfix     | local         | production_commit |
      | production | local, remote | production_commit |
