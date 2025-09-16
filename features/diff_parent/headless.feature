@smoke
Feature: view changes in headless state

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | branch-1 | feature | main   | local     |
    And the commits
      | BRANCH   | LOCATION | MESSAGE   |
      | branch-1 | local    | commit 1a |
      |          |          | commit 1b |
    And the current branch is "branch-1"
    And I ran "git checkout HEAD^"
    When I run "git-town diff-parent"

  @this
  Scenario: feature branch
    Then Git Town runs the commands
      | BRANCH   | COMMAND                             |
      | branch-1 | git diff --merge-base main branch-1 |
