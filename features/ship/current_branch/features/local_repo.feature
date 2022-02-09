Feature: ship a feature branch in a local repo

  Background:
    Given a feature branch "feature"
    And my repo does not have an origin
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And I am on the "feature" branch
    When I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                      |
      | feature | git merge --no-edit main     |
      |         | git checkout main            |
      | main    | git merge --squash feature   |
      |         | git commit -m "feature done" |
      |         | git branch -D feature        |
    And I am now on the "main" branch
    And the branches are now
      | REPOSITORY | BRANCHES |
      | local      | main     |
    And now these commits exist
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | feature done |
    And Git Town is now aware of no branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
      |        | git revert {{ sha 'feature done' }}           |
      |        | git checkout feature                          |
    And I am now on the "feature" branch
    And now these commits exist
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | feature done          |
      |         |          | Revert "feature done" |
      | feature | local    | feature commit        |
    And the initial branches and hierarchy exist
