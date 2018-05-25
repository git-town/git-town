Feature: git town-hack: starting a new feature from the main branch (with remote repo)

  As a developer working on a new feature on the main branch
  I want to be able to create a new up-to-date feature branch and continue my work there
  So that my work can exist on its own branch, code reviews remain effective, and my team productive.


  Background:
    Given my repository has the perennial branches "production"
    And the following commit exists in my repository
      | BRANCH     | LOCATION | MESSAGE           |
      | main       | remote   | main_commit       |
      | production | remote   | production_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file


  Scenario: selecting the default branch (the main development branch)
    When I run `git-town hack -p new-feature` and answer the prompts:
      | PROMPT                                            | ANSWER  |
      | Please specify the parent branch of 'new-feature' | [ENTER] |
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git fetch --prune           |
      |             | git add -A                  |
      |             | git stash                   |
      |             | git rebase origin/main      |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I end up on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repository has the following commits
      | BRANCH      | LOCATION         | MESSAGE           |
      | main        | local and remote | main_commit       |
      | new-feature | local            | main_commit       |
      | production  | remote           | production_commit |


    Scenario: selecting another branch
      When I run `git-town hack -p hotfix` and answer the prompts:
        | PROMPT                                       | ANSWER        |
        | Please specify the parent branch of 'hotfix' | [DOWN][ENTER] |
      Then it runs the commands
        | BRANCH     | COMMAND                      |
        | main       | git fetch --prune            |
        |            | git add -A                   |
        |            | git stash                    |
        |            | git checkout production      |
        | production | git rebase origin/production |
        |            | git branch hotfix production |
        |            | git checkout hotfix          |
        | hotfix     | git stash pop                |
      And I end up on the "hotfix" branch
      And my workspace still contains my uncommitted file
      And my repository has the following commits
        | BRANCH     | LOCATION         | MESSAGE           |
        | main       | remote           | main_commit       |
        | hotfix     | local            | production_commit |
        | production | local and remote | production_commit |
