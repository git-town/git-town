Feature: git sync --all: handling merge conflicts between feature branch and main branch with open changes

  Background:
    Given I have feature branches named "feature1" and "feature2"
    And the following commits exist in my repository
      | BRANCH   | LOCATION         | MESSAGE                | FILE NAME            | FILE CONTENT            |
      | main     | remote           | main commit            | conflicting_file     | main content            |
      | feature1 | local and remote | feature1 commit        | feature1_file        | feature1 content        |
      | feature2 | local            | feature2 local commit  | conflicting_file     | feature2 local content  |
      |          | remote           | feature2 remote commit | feature2_remote_file | feature2 remote content |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `git sync --all`


  Scenario: result
    Then it runs the Git commands
      | BRANCH   | COMMAND                             |
      | main     | git fetch --prune                   |
      |          | git stash -u                        |
      |          | git rebase origin/main              |
      |          | git checkout feature1               |
      | feature1 | git merge --no-edit origin/feature1 |
      |          | git merge --no-edit main            |
      |          | git push                            |
      |          | git checkout feature2               |
      | feature2 | git merge --no-edit origin/feature2 |
      |          | git merge --no-edit main            |
    And I get the error
      """
      To abort, run "git sync --abort".
      To continue after you have resolved the conflicts, run "git sync --continue".
      To skip the sync of the 'feature2' branch, run "git sync --skip".
      """
    And I end up on the "feature2" branch
    And my uncommitted file is still stashed away
    And my repo has a merge in progress


  Scenario: aborting
    When I run `git sync --abort`
    Then it runs the Git commands
      | BRANCH   | COMMAND                                             |
      | feature2 | git merge --abort                                   |
      |          | git reset --hard <%= sha 'feature2 local commit' %> |
      |          | git checkout feature1                               |
      | feature1 | git checkout main                                   |
      | main     | git stash pop                                       |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                           | FILE NAME            |
      | main     | local and remote | main commit                       | conflicting_file     |
      | feature1 | local and remote | feature1 commit                   | feature1_file        |
      |          |                  | main commit                       | conflicting_file     |
      |          |                  | Merge branch 'main' into feature1 |                      |
      | feature2 | local            | feature2 local commit             | conflicting_file     |
      |          | remote           | feature2 remote commit            | feature2_remote_file |


  Scenario: skipping
    When I run `git sync --skip`
    Then it runs the Git commands
      | BRANCH   | COMMAND                                             |
      | feature2 | git merge --abort                                   |
      |          | git reset --hard <%= sha 'feature2 local commit' %> |
      |          | git checkout main                                   |
      | main     | git stash pop                                       |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                           | FILE NAME            |
      | main     | local and remote | main commit                       | conflicting_file     |
      | feature1 | local and remote | feature1 commit                   | feature1_file        |
      |          |                  | main commit                       | conflicting_file     |
      |          |                  | Merge branch 'main' into feature1 |                      |
      | feature2 | local            | feature2 local commit             | conflicting_file     |
      |          | remote           | feature2 remote commit            | feature2_remote_file |


  Scenario: continuing without resolving the conflicts
    When I run `git sync --continue`
    Then it runs no Git commands
    And I get the error "You must resolve the conflicts before continuing the git sync"
    And I am still on the "feature2" branch
    And my uncommitted file is still stashed away
    And my repo still has a merge in progress


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run `git sync --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND              |
      | feature2 | git commit --no-edit |
      |          | git push             |
      |          | git checkout main    |
      | main     | git stash pop        |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                                                      | FILE NAME            |
      | main     | local and remote | main commit                                                  | conflicting_file     |
      | feature1 | local and remote | feature1 commit                                              | feature1_file        |
      |          |                  | main commit                                                  | conflicting_file     |
      |          |                  | Merge branch 'main' into feature1                            |                      |
      | feature2 | local and remote | feature2 local commit                                        | conflicting_file     |
      |          |                  | feature2 remote commit                                       | feature2_remote_file |
      |          |                  | Merge remote-tracking branch 'origin/feature2' into feature2 |                      |
      |          |                  | main commit                                                  | conflicting_file     |
      |          |                  | Merge branch 'main' into feature2                            |                      |



  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    And I run `git commit --no-edit; git sync --continue`
    Then it runs the Git commands
      | BRANCH   | COMMAND           |
      | feature2 | git push          |
      |          | git checkout main |
      | main     | git stash pop     |
    And I end up on the "main" branch
    And I again have my uncommitted file
    And I have the following commits
      | BRANCH   | LOCATION         | MESSAGE                                                      | FILE NAME            |
      | main     | local and remote | main commit                                                  | conflicting_file     |
      | feature1 | local and remote | feature1 commit                                              | feature1_file        |
      |          |                  | main commit                                                  | conflicting_file     |
      |          |                  | Merge branch 'main' into feature1                            |                      |
      | feature2 | local and remote | feature2 local commit                                        | conflicting_file     |
      |          |                  | feature2 remote commit                                       | feature2_remote_file |
      |          |                  | Merge remote-tracking branch 'origin/feature2' into feature2 |                      |
      |          |                  | main commit                                                  | conflicting_file     |
      |          |                  | Merge branch 'main' into feature2                            |                      |
