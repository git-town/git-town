Feature: ignore files

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE   | FILE NAME  | FILE CONTENT |
      | feature | local    | my commit | .gitignore | ignored      |
    And the current branch is "feature"
    And an uncommitted file "test/ignored/important" with content "changed ignored file"
    When I run "git-town sync"

  Scenario: result
    Then file "test/ignored/important" still has content "changed ignored file"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                               |
      | feature | git push --force-with-lease origin {{ sha 'initial commit' }}:feature |
    And the initial branches and lineage exist now
    And the initial commits exist now
