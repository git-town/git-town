Feature: handle uncommitted files

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS |
      | branch-1 | feature | main     | local     |
      | branch-2 | feature | branch-1 | local     |
    And the current branch is "branch-2"
    And an uncommitted file
    When I run "git-town walk --stack echo hello"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | branch-2 | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git checkout branch-1       |
      | branch-1 | echo hello                  |
      |          | git checkout branch-2       |
      | branch-2 | echo hello                  |
      |          | git stash pop               |
      |          | git restore --staged .      |
    And Git Town prints:
      """
      Branch walk done.
      """
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                     |
      | branch-2 | git add -A                  |
      |          | git stash -m "Git Town WIP" |
      |          | git stash pop               |
      |          | git restore --staged .      |
