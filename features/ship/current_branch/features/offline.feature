Feature: offline mode

  Background:
    Given offline mode is enabled
    And a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    When I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git branch -D feature              |
    And the current branch is now "main"
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE        |
      | main    | local    | feature done   |
      | feature | origin   | feature commit |
    And Git Town is now aware of no branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git reset --hard {{ sha 'Initial commit' }}   |
      |         | git checkout feature                          |
    And the current branch is now "feature"
    And now the initial commits exist
    And the initial branches and hierarchy exist
