Feature: delete a remote only branch

  Background:
    Given a remote feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | origin   | feature commit |
    And the current branch is "main"
    And I run "git-town sync"
    When I run "git-town kill feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git push origin :feature |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And Git Town is now aware of no branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                 |
      | main   | git push origin {{ sha-in-origin 'feature commit' }}:refs/heads/feature |
    And the initial branches exist
    And Git Town is still aware of no branch hierarchy
