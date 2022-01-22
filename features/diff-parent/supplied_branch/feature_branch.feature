Feature: git town-diff-parent: diffing a given feature branch

  To know whether my global branch setup is correct
  When working with nested feature branches
  I want to see the changes that a particular feature branch makes.

  Scenario:
    Given my repo has a feature branch named "feature-1"
    And my repo has a feature branch named "feature-2" as a child of "feature-1"
    When I run "git-town diff-parent feature-2"
    Then it runs the commands
      | BRANCH | COMMAND                       |
      | main   | git diff feature-1..feature-2 |
