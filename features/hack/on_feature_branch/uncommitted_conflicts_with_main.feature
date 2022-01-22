Feature: git town-hack: resolving conflicts between uncommitted changes and the main branch

  Background:
    Given my repo has a feature branch named "feature1"
    And the following commits exist in my repo
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, remote | conflicting commit | conflicting_file | main content |
    And I am on the "feature1" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town hack feature2"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                     |
      | feature1 | git fetch --prune --tags    |
      |          | git add -A                  |
      |          | git stash                   |
      |          | git checkout main           |
      | main     | git rebase origin/main      |
      |          | git branch feature2 main |
      |          | git checkout feature2    |
      | feature2 | git stash pop               |
    And it prints the error:
      """
      conflicts between your uncommmitted changes and the main branch
      """
    And the file "conflicting_file" contains unresolved conflicts

  # TODO: abort

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue" and close the editor
    Then it runs no commands
    And I am now on the "feature2" branch
    And my workspace now contains the file "conflicting_file" with content "resolved content"
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE            | FILE NAME        |
      | main        | local, remote | conflicting commit | conflicting_file |
      | feature2 | local         | conflicting commit | conflicting_file |
