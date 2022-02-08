Feature: dry run

  Background:
    Given my repo has a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | remote main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | remote feature commit |
    And I am on the "feature" branch
    When I run "git-town sync --dry-run"

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
      |         | git push                           |
    And I am still on the "feature" branch
    And now the initial commits exist
