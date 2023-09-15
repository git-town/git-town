Feature: undoing an offline ship with additional commits to main

  Background:
    Given offline mode is enabled
    And the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    When I run "git-town ship -m 'feature done'"
    And I add commit "additional commit" to the "main" branch

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
    And it prints the error:
      """
      cannot reset branch "main"
      """
    And it prints the error:
      """
      it received additional commits in the meantime
      """
    And the current branch is now "main"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local         | feature done          |
      |         |               | additional commit     |
      |         |               | Revert "feature done" |
      | feature | local, origin | feature commit        |
    And the initial branches and hierarchy exist
