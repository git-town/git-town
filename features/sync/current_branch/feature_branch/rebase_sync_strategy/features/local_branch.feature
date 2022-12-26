Feature: sync the current feature branch without a tracking branch

  Background:
    Given setting "sync-strategy" is "rebase"
    And the current branch is a local feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      |         | origin   | origin main commit   |
      | feature | local    | local feature commit |
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git push                   |
      |         | git checkout feature       |
      | feature | git rebase main            |
      |         | git push -u origin feature |
    And all branches are now synchronized
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE              |
      | main    | local, origin | origin main commit   |
      |         |               | local main commit    |
      | feature | local, origin | origin main commit   |
      |         |               | local main commit    |
      |         |               | local feature commit |
    And the branches are now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
