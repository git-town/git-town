Feature: offline mode

  Background:
    Given offline mode is enabled
    And a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |
    And the current branch is "feature"
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE                          |
      | main    | local    | local main commit                |
      |         | origin   | origin main commit               |
      | feature | local    | local feature commit             |
      |         |          | local main commit                |
      |         |          | Merge branch 'main' into feature |
      |         | origin   | origin feature commit            |
