Feature: ignore files

  Background:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
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
      | BRANCH  | COMMAND                                                               |
      | feature | git push --force-with-lease origin {{ sha 'initial commit' }}:feature |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
