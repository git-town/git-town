Feature: handling merge conflicts between feature branch and main branch

  Background:
    Given my repo has the feature branches "feature-1", "feature-2", and "feature-3"
    And my repo contains the commits
      | BRANCH    | LOCATION      | MESSAGE                 | FILE NAME            | FILE CONTENT             |
      | main      | remote        | main commit             | conflicting_file     | main content             |
      | feature-1 | local, remote | feature-1 commit        | feature1_file        | feature-1 content        |
      | feature-2 | local         | feature-2 local commit  | conflicting_file     | feature-2 local content  |
      |           | remote        | feature-2 remote commit | feature2_remote_file | feature-2 remote content |
      | feature-3 | remote        | feature-3 commit        | feature3_file        | feature-3 content        |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                              |
      | main      | git fetch --prune --tags             |
      |           | git add -A                           |
      |           | git stash                            |
      |           | git rebase origin/main               |
      |           | git checkout feature-1               |
      | feature-1 | git merge --no-edit origin/feature-1 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git checkout feature-2               |
      | feature-2 | git merge --no-edit origin/feature-2 |
      |           | git merge --no-edit main             |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      To continue by skipping the current branch, run "git-town skip".
      """
    And I am now on the "feature-2" branch
    And my uncommitted file is stashed
    And my repo now has a merge in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH    | COMMAND                                             |
      | feature-2 | git merge --abort                                   |
      |           | git reset --hard {{ sha 'feature-2 local commit' }} |
      |           | git checkout feature-1                              |
      | feature-1 | git checkout main                                   |
      | main      | git stash pop                                       |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH    | LOCATION      | MESSAGE                            |
      | main      | local, remote | main commit                        |
      | feature-1 | local, remote | feature-1 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-1 |
      | feature-2 | local         | feature-2 local commit             |
      |           | remote        | feature-2 remote commit            |
      | feature-3 | remote        | feature-3 commit                   |
    And my repo still has the following committed files
      | BRANCH    | NAME             | CONTENT                 |
      | main      | conflicting_file | main content            |
      | feature-1 | conflicting_file | main content            |
      |           | feature1_file    | feature-1 content       |
      | feature-2 | conflicting_file | feature-2 local content |

  Scenario: skipping
    When I run "git-town skip"
    Then it runs the commands
      | BRANCH    | COMMAND                                             |
      | feature-2 | git merge --abort                                   |
      |           | git reset --hard {{ sha 'feature-2 local commit' }} |
      |           | git checkout feature-3                              |
      | feature-3 | git merge --no-edit origin/feature-3                |
      |           | git merge --no-edit main                            |
      |           | git push                                            |
      |           | git checkout main                                   |
      | main      | git push --tags                                     |
      |           | git stash pop                                       |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH    | LOCATION      | MESSAGE                            |
      | main      | local, remote | main commit                        |
      | feature-1 | local, remote | feature-1 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-1 |
      | feature-2 | local         | feature-2 local commit             |
      |           | remote        | feature-2 remote commit            |
      | feature-3 | local, remote | feature-3 commit                   |
      |           |               | main commit                        |
      |           |               | Merge branch 'main' into feature-3 |
    And my repo now has the following committed files
      | BRANCH    | NAME             | CONTENT                 |
      | main      | conflicting_file | main content            |
      | feature-1 | conflicting_file | main content            |
      |           | feature1_file    | feature-1 content       |
      | feature-2 | conflicting_file | feature-2 local content |
      | feature-3 | conflicting_file | main content            |
      |           | feature3_file    | feature-3 content       |

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And I am still on the "feature-2" branch
    And my uncommitted file is stashed
    And my repo still has a merge in progress

  Scenario: continuing after resolving the conflicts
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH    | COMMAND                              |
      | feature-2 | git commit --no-edit                 |
      |           | git push                             |
      |           | git checkout feature-3               |
      | feature-3 | git merge --no-edit origin/feature-3 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git checkout main                    |
      | main      | git push --tags                      |
      |           | git stash pop                        |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And all branches are now synchronized
    And my repo now has the following committed files
      | BRANCH    | NAME                 | CONTENT                  |
      | main      | conflicting_file     | main content             |
      | feature-1 | conflicting_file     | main content             |
      |           | feature1_file        | feature-1 content        |
      | feature-2 | conflicting_file     | resolved content         |
      |           | feature2_remote_file | feature-2 remote content |
      | feature-3 | conflicting_file     | main content             |
      |           | feature3_file        | feature-3 content        |

  Scenario: continuing after resolving the conflicts and committing
    When I resolve the conflict in "conflicting_file"
    And I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH    | COMMAND                              |
      | feature-2 | git push                             |
      |           | git checkout feature-3               |
      | feature-3 | git merge --no-edit origin/feature-3 |
      |           | git merge --no-edit main             |
      |           | git push                             |
      |           | git checkout main                    |
      | main      | git push --tags                      |
      |           | git stash pop                        |
    And I am now on the "main" branch
    And all branches are now synchronized
    And my workspace has the uncommitted file again
