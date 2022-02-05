Feature: must provide a commit message

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And I am on the "feature" branch
    When I run "git-town ship" and close the editor

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit                         |
      |         | git reset --hard                   |
      |         | git checkout feature               |
      | feature | git checkout main                  |
      | main    | git checkout feature               |
    And it prints the error:
      """
      aborted because commit exited with error
      """
    And I am still on the "feature" branch
    And my repo is left with my initial commits
    And Git Town still has the initial branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature" branch
    And my repo is left with my initial commits
    And Git Town still has the initial branch hierarchy
