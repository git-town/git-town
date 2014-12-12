Feature: git prune-branches: leaves used feature branches when called on the main branch (without open changes)

  As a developer having feature branches with commits
  I want them all to survive a "prune branch" command
  So that I can keep my repository clean without loosing work.


  Background:
    Given I have a feature branch named "my-feature" ahead of main
    And my coworker has a feature branch named "charlies-feature" ahead of main
    And I am on the "main" branch
    When I run `git prune-branches`


  Scenario: result
    Then it runs the Git commands
      | BRANCH | COMMAND           |
      | main   | git fetch --prune |
    Then I end up on the "main" branch
    And the existing branches are
      | REPOSITORY | BRANCHES                           |
      | local      | main, my-feature                   |
      | remote     | main, my-feature, charlies-feature |
      | coworker   | main, charlies-feature             |
