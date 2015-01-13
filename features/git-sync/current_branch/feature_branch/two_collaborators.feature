Feature: git sync: collaborative feature branch syncing

  As a developer collaborating with others on a feature
  I want each person to be able to sync their changes with the rest of the team
  So that our collaboration is effective.


  Background:
    Given I have a feature branch named "feature"
    And my coworker fetches updates
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE   | FILE NAME |
      | feature | local    | my commit | my_file   |
    And the following commits exist in my coworker's repository
      | BRANCH  | LOCATION | MESSAGE         | FILE NAME     |
      | feature | local    | coworker commit | coworker_file |

  Scenario: result
    And I am on the "feature" branch
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE   | FILE NAME |
      | feature | local and remote | my commit | my_file   |

    Given my coworker is on the "feature" branch
    When my coworker runs `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
      | feature | git push                           |
    And now my coworker has the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME     |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |               |
      | feature |                  | my commit                                                  | my_file       |
      | feature |                  | coworker commit                                            | coworker_file |

    Given I am on the "feature" branch
    When I run `git sync`
    Then it runs the Git commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune                  |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      | main    | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      | feature | git merge --no-edit main           |
    And now I have the following commits
      | BRANCH  | LOCATION         | MESSAGE                                                    | FILE NAME     |
      | feature | local and remote | Merge remote-tracking branch 'origin/feature' into feature |               |
      | feature |                  | my commit                                                  | my_file       |
      | feature |                  | coworker commit                                            | coworker_file |
