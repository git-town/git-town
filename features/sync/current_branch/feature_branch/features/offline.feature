Feature: offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | remote   | remote main commit    |
      | feature | local    | local feature commit  |
      |         | remote   | remote feature commit |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | local main commit                |
      |         | remote   | remote main commit               |
      | feature | local    | local feature commit             |
      |         |          | local main commit                |
      |         |          | Merge branch 'main' into feature |
      |         | remote   | remote feature commit            |
