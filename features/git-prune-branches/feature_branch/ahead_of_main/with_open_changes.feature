Feature: git prune-branches: keeps used feature branches when run on a feature branch without open changes

  As a developer having feature branches with commits
  I want them all to survive a "prune branch" command
  So that I can keep my repository clean without loosing work.


  Background:
    Given I have a feature branch named "feature" ahead of main
    And I have a feature branch named "stale_feature" behind main
    And I am on the "feature" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                        |
      | feature | git fetch --prune              |
      | feature | git stash -u                   |
      | feature | git checkout main              |
      | main    | git push origin :stale_feature |
      | main    | git branch -d stale_feature    |
      | main    | git checkout feature           |
      | feature | git stash pop                  |
    Then I end up on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | remote     | main, feature |
      | coworker   | main          |


  Scenario: [note]
    Given TODO: this feature specifies two behaviors: that branches ahead of main are not deleted, and that old branches are. Specify only one feature per spec. Here we should have two full feature branches, which both should survive the cleanup.

  Scenario: undo
    Given TODO: the user should be able to undo this command, or if there is no undo we should make sure he can't run undo here.
