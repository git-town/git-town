@skipWindows
Feature: aborting the ship by empty commit message

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And I am on the "feature" branch
    When I run "git-town ship" and enter an empty commit message

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
    And my repo is left with my original commits
    And Git Town is still aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature" branch
    And my repo is left with my original commits
    And Git Town is still aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
