Feature: prototype another remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME           | TYPE   | PARENT | LOCATIONS |
      | remote-feature | (none) | main   | origin    |
    And I run "git fetch"
    When I run "git-town prototype remote-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                     |
      | main   | git checkout remote-feature |
    And it prints:
      """
      branch "remote-feature" is now a prototype branch
      """
    And the current branch is now "remote-feature"
    And branch "remote-feature" is now prototype

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH         | COMMAND                      |
      | remote-feature | git checkout main            |
      | main           | git branch -D remote-feature |
    And the current branch is now "main"
    And there are now no observed branches
    And the initial branches exist now
