Feature: ignore files

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE   | FILE NAME  | FILE CONTENT |
      | feature | local    | my commit | .gitignore | ignored      |
    And an uncommitted file with name "test/ignored/important" and content "changed ignored file"
    When I run "git-town sync"

  Scenario: result
    Then file "test/ignored/important" still has content "changed ignored file"

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND              |
      | feature | git checkout main    |
      | main    | git checkout feature |
    And the current branch is still "feature"
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE   |
      | feature | local, origin | my commit |
    And the initial branches and hierarchy exist
