Feature: resolving conflicts between the supplied feature branch and the main branch

  Background:
    Given my repo has the feature branches "feature" and "other-feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file
    And I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune --tags           |
      |               | git add -A                         |
      |               | git stash                          |
      |               | git checkout main                  |
      | main          | git rebase origin/main             |
      |               | git push                           |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      |               | git merge --no-edit main           |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And I am now on the "feature" branch
    And my uncommitted file is stashed
    And my repo now has a merge in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH        | COMMAND                    |
      | feature       | git merge --abort          |
      |               | git checkout main          |
      | main          | git checkout other-feature |
      | other-feature | git stash pop              |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And there is no merge in progress
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        |
      | main    | local, remote | conflicting main commit    | conflicting_file |
      | feature | local         | conflicting feature commit | conflicting_file |

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH        | COMMAND                      |
      | feature       | git commit --no-edit         |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other-feature   |
      | other-feature | git stash pop                |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE                 | FILE NAME        |
      | main   | local, remote | conflicting main commit | conflicting_file |
      |        |               | feature done            | conflicting_file |

  Scenario: continuing after resolving the conflicts and comitting
    Given I resolve the conflict in "conflicting_file"
    When I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH        | COMMAND                      |
      | feature       | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git push                     |
      |               | git push origin :feature     |
      |               | git branch -D feature        |
      |               | git checkout other-feature   |
      | other-feature | git stash pop                |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
      | remote     | main, other-feature |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE                 | FILE NAME        |
      | main   | local, remote | conflicting main commit | conflicting_file |
      |        |               | feature done            | conflicting_file |
