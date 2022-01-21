Feature: git-sync: on a feature branch in a repository with a submodule

  As a developer of a repository that uses submodules
  I want that "git sync" ignores uncommitted changes in the submodules
  So that I can focus on changes in my codebase and I can deal with the submodules separately.

  Background:
    Given my repo has a submodule
    And my repo has a feature branch named "feature"
    And I am on the "feature" branch
    And my workspace has an uncommitted file with name "submodule/changed" and content "content"
    And inspect the repo
    When I run "git-town sync"


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git add -A                         |
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
      | BRANCH | LOCATION | MESSAGE |
