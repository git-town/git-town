Feature: handle conflicts between the shipped branch and the main branch

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I am on the "feature" branch
    And I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And I am still on the "feature" branch
    And my repo now has a merge in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git merge --abort    |
      |         | git checkout main    |
      | main    | git checkout feature |
    And I am still on the "feature" branch
    And there is no merge in progress
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                    | FILE NAME        | FILE CONTENT    |
      | main    | local, remote | conflicting main commit    | conflicting_file | main content    |
      | feature | local         | conflicting feature commit | conflicting_file | feature content |
    And Git Town is still aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git commit --no-edit         |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES |
      | local      | main     |
      | remote     | main     |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE                 | FILE NAME        | FILE CONTENT     |
      | main   | local, remote | conflicting main commit | conflicting_file | main content     |
      |        |               | feature done            | conflicting_file | resolved content |
    And Git Town now has no branch hierarchy information

  Scenario: continuing after resolving the conflicts and committing
    Given I resolve the conflict in "conflicting_file"
    When I run "git commit --no-edit"
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git push                     |
      |         | git push origin :feature     |
      |         | git branch -D feature        |
    And I am now on the "main" branch
