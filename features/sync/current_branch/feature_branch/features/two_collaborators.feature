Feature: collaborative feature branch syncing

  Scenario:
    Given I am collaborating with a coworker
    And my repo has a feature branch "feature"
    And my coworker fetches updates
    And my coworker sets the parent branch of "feature" as "main"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE         |
      | feature | local    | my commit       |
      |         | coworker | coworker commit |
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
    And my repo now has the commits
      | BRANCH  | LOCATION      | MESSAGE         |
      | feature | local, remote | my commit       |
      |         | coworker      | coworker commit |

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
    And my repo now has the commits
      | BRANCH  | LOCATION                | MESSAGE                                                    |
      | feature | local, coworker, remote | my commit                                                  |
      |         | coworker, remote        | coworker commit                                            |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature |

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
    And my repo now has the commits
      | BRANCH  | LOCATION                | MESSAGE                                                    |
      | feature | local, coworker, remote | coworker commit                                            |
      |         |                         | my commit                                                  |
      |         |                         | Merge remote-tracking branch 'origin/feature' into feature |
