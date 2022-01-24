Feature: git-town sync: collaborative feature branch syncing

  Background:
    Given I am collaborating with a coworker
    And my repo has a feature branch named "feature"
    And my coworker fetches updates
    And my coworker sets the parent branch of "feature" as "main"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME     |
      | feature | local    | my commit       | my_file       |
      |         | coworker | coworker commit | coworker_file |

  Scenario: result
    And I am on the "feature" branch
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE         | FILE NAME     |
      | feature | local, remote | my commit       | my_file       |
      |         | coworker      | coworker commit | coworker_file |

    Given my coworker is on the "feature" branch
    When my coworker runs "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
    And my repo now has the following commits
      | BRANCH  | LOCATION                | MESSAGE                                                    | FILE NAME     |
      | feature | local, coworker, remote | my commit                                                  | my_file       |
      |         | coworker, remote        | coworker commit                                            | coworker_file |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature |               |

    Given I am on the "feature" branch
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And my repo now has the following commits
      | BRANCH  | LOCATION                | MESSAGE                                                    | FILE NAME     |
      | feature | local, coworker, remote | coworker commit                                            | coworker_file |
      |         |                         | my commit                                                  | my_file       |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature |               |
