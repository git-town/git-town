Feature: git hack: don't push branch to remote upon creation

  As a developer starting work on a new private feature
  I don't want git-hack to push my branch to the remote repo
  So that my CI server doesn't unnecessarily run tests for commits without changes


  Background:
    Given my repo has an upstream repo
    And my repository has the "hack-push-flag" configuration set to "local"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    When I run `git hack private-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                              |
      | main   | git fetch --prune                    |
      |        | git rebase origin/main               |
      |        | git fetch upstream                   |
      |        | git rebase upstream/main             |
      |        | git push                             |
      |        | git checkout -b private-feature main |
    And I am still on the "private-feature" branch
    And I have the following commits
      | BRANCH          | LOCATION                    | MESSAGE         |
      | main            | local, remote, and upstream | upstream commit |
      | private-feature | local                       | upstream commit |
