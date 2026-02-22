@skipWindows
Feature: no TTY

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
    And the current branch is "branch-2"
    And an uncommitted file "changes" with content "my changes"
    And I ran "git add changes"
    When I run "git-town commit --down -m commit-1b" in a non-TTY shell

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                           |
      | branch-2 | git checkout branch-1             |
      | branch-1 | git commit -m commit-1b           |
      |          | git checkout branch-2             |
      | branch-2 | git merge --no-edit --ff branch-1 |
