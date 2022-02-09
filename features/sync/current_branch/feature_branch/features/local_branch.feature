Feature: sync the current feature branch without a tracking branch

  Background:
    Given a local feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE              |
      | main    | local    | local main commit    |
      |         | origin   | origin main commit   |
      | feature | local    | local feature commit |
    And I am on the "feature" branch
    When I run "git-town sync"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                    |
      | feature | git fetch --prune --tags   |
      |         | git checkout main          |
      | main    | git rebase origin/main     |
      |         | git push                   |
      |         | git checkout feature       |
      | feature | git merge --no-edit main   |
      |         | git push -u origin feature |
    And all branches are now synchronized
    And I am still on the "feature" branch
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                          |
      | main    | local, origin | origin main commit               |
      |         |               | local main commit                |
      | feature | local, origin | local feature commit             |
      |         |               | origin main commit               |
      |         |               | local main commit                |
      |         |               | Merge branch 'main' into feature |
    And the existing branches are
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
