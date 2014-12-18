Feature: git prune-branches: leave used feature branches when called on the main branch (with open changes)

  As a developer having feature branches with commits
  I want them all to survive a "prune branch" command
  So that I can keep my repository clean without loosing work.


  Background:
    Given I have a feature branch named "my-feature" ahead of main
    And my coworker has a feature branch named "co-feature" ahead of main
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
      | REPOSITORY | BRANCHES                     |
      | local      | main, my-feature             |
      | remote     | main, my-feature, co-feature |
      | coworker   | main, co-feature             |
