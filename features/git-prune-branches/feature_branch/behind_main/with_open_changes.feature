Feature: git prune-branches: removes stale feature branches when run on a feature branch with open changes

  As a developer having empty feature branches
  I want them all to be cleaned out
  So that all my remaining branches are relevant and I can focus on my current work.


  Background:
    Given I have a feature branch named "feature" behind main
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
    Given TODO: this feature specifies two behaviors: that branches ahead of main are not deleted, and that old branches are. Specify only one feature per spec.
