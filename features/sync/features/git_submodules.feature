Feature: on a feature branch in a repository with a submodule that has uncommitted changes

  Background:
    Given my repo has a submodule
    And my repo has a feature branch "feature"
    And I am on the "feature" branch
    And my workspace has an uncommitted file with name "submodule/file" and content "a change in the submodule"
    When I run "git-town sync"

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
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file
