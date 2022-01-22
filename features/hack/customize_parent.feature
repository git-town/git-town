@skipWindows
Feature: git town hack: customize the parent branch

  To allow hotfixes for other branches than the main branch
  When creating a new feature branch
  I want to be able to customize the parent for the new branch.

  Scenario:
    Given my repo has the perennial branch "production"
    And the following commits exist in my repo
      | BRANCH     | LOCATION | MESSAGE           |
      | production | remote   | production_commit |
    And I am on the "main" branch
    When I run "git-town hack -p hotfix" and answer the prompts:
      | PROMPT                                       | ANSWER        |
      | Please specify the parent branch of 'hotfix' | [DOWN][ENTER] |
    Then it runs the commands
      | BRANCH     | COMMAND                      |
      | main       | git fetch --prune --tags     |
      |            | git checkout production      |
      | production | git rebase origin/production |
      |            | git branch hotfix production |
      |            | git checkout hotfix          |
    And I am now on the "hotfix" branch
    And my repo now has the following commits
      | BRANCH     | LOCATION      | MESSAGE           |
      | hotfix     | local         | production_commit |
      | production | local, remote | production_commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT     |
      | hotfix | production |
