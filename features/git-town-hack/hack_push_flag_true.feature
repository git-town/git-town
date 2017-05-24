Feature: git town-hack: push branch to remote upon creation

  As a developer starting work on a new private feature
  I want to be able to configure git-hack to push my branch to the remote repo
  So that ...


  Background:
    Given my repository has the "hack-push-flag" configuration set to "true"
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE       |
      | main   | remote   | remote commit |
    And I am on the "main" branch
    When I run `git-town hack feature`


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | main    | git fetch --prune            |
      |         | git rebase origin/main       |
      |         | git checkout -b feature main |
      | feature | git push -u origin feature   |
    And I end up on the "feature" branch
    And I have the following commits
      | BRANCH  | LOCATION         | MESSAGE       |
      | main    | local and remote | remote commit |
      | feature | local and remote | remote commit |
