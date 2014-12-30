Feature: git-sync-all from the main branch

  Background:
    Given I have feature branches named "feature1" and "feature2"
    And the following commits exist in my repository
      | branch   | location         | message         | file name     |
      | main     | remote           | main commit     | main_file     |
      | feature1 | local and remote | feature1 commit | feature1_file |
      | feature2 | local and remote | feature2 commit | feature2_file |
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                             |
      | main     | git fetch --prune                   |
      | main     | git stash -u                        |
      | main     | git rebase origin/main              |
      | main     | git checkout feature1               |
      | feature1 | git merge --no-edit origin/feature1 |
      | feature1 | git merge --no-edit main            |
      | feature1 | git push                            |
      | feature1 | git checkout feature2               |
      | feature2 | git merge --no-edit origin/feature2 |
      | feature2 | git merge --no-edit main            |
      | feature2 | git push                            |
      | feature2 | git checkout main                   |
      | main     | git stash pop                       |
    And I am still on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And all branches are now synchronized
    And I have the following commits
      | branch   | location         | message                           | files         |
      | main     | local and remote | main commit                       | main_file     |
      | feature1 | local and remote | Merge branch 'main' into feature1 |               |
      |          |                  | main commit                       | main_file     |
      |          |                  | feature1 commit                   | feature1_file |
      | feature2 | local and remote | Merge branch 'main' into feature2 |               |
      |          |                  | main commit                       | main_file     |
      |          |                  | feature2 commit                   | feature2_file |
