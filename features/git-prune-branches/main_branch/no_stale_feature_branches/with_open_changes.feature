Feature: git prune-branches: leave used feature branches when called on the main branch (with open changes)

  As a developer having feature branches with commits
  I want them all to survive a "prune branch" command
  So that I can keep my repository clean without loosing work.


  Background:
    Given I have a feature branch named "feature1" ahead of main
    And my coworker has a feature branch named "feature2" ahead of main
    And I am on the "main" branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    Then I end up on the "main" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the existing branches are
      | REPOSITORY | BRANCHES                 |
      | local      | main, feature1           |
      | remote     | main, feature1, feature2 |
      | coworker   | main, feature2           |
