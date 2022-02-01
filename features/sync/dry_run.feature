Feature: dry run

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | remote   | remote main commit    |
      | feature | local    | local feature commit  |
      |         | remote   | remote feature commit |
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
    And my repo is left with my original commits
    And Git Town is still aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
