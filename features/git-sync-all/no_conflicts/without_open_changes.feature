Feature: git-sync-all

  Background:
    Given I have branches named "feature1", "feature2", and "production"
    And my non-feature branch is "production"
    And the following commits exist in my repository
      | branch     | location         | message                  | file name              |
      | main       | remote           | main commit              | main_file              |
      | feature1   | local and remote | feature1 commit          | feature1_file          |
      | feature2   | local and remote | feature2 commit          | feature2_file          |
      | production | local            | production local commit  | production_local_file  |
      | production | remote           | production remote commit | production_remote_file |
    And I am on the "main" branch
    When I run `git sync-all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH     | COMMAND                             |
      | main       | git fetch --prune                   |
      | main       | git rebase origin/main              |
      | main       | git checkout feature1               |
      | feature1   | git merge --no-edit origin/feature1 |
      | feature1   | git merge --no-edit main            |
      | feature1   | git push                            |
      | feature1   | git checkout feature2               |
      | feature2   | git merge --no-edit origin/feature2 |
      | feature2   | git merge --no-edit main            |
      | feature2   | git push                            |
      | feature2   | git checkout production             |
      | production | git rebase origin/production        |
      | production | git push                            |
      | production | git checkout main                   |
    And I am still on the "main" branch
    And all branches are now synchronized
    And I have the following commits
      | branch     | location         | message                           | files                  |
      | main       | local and remote | main commit                       | main_file              |
      | feature1   | local and remote | Merge branch 'main' into feature1 |                        |
      |            |                  | main commit                       | main_file              |
      |            |                  | feature1 commit                   | feature1_file          |
      | feature2   | local and remote | Merge branch 'main' into feature2 |                        |
      |            |                  | main commit                       | main_file              |
      |            |                  | feature2 commit                   | feature2_file          |
      | production | local and remote | production local commit           | production_local_file  |
      | production | local and remote | production remote commit          | production_remote_file |