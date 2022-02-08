Feature: offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | remote main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | remote feature commit |
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And I am still on the "feature" branch
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | local main commit                |
      |         | origin   | remote main commit               |
      | feature | local    | local feature commit             |
      |         |          | local main commit                |
      |         |          | Merge branch 'main' into feature |
      |         | origin   | remote feature commit            |
