Feature: preserve uncommitted files

  Scenario:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE                      |
      | main    | local    | local main commit            |
      |         | remote   | remote main commit           |
      | feature | local    | local parent feature commit  |
      |         | remote   | remote parent feature commit |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git add -A                         |
      |         | git stash                          |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git stash pop                      |
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
