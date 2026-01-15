@skipWindows
Feature: open the page of an already existing proposal

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | URL                                           |
      |  1 | feature       | main          | https://github.com/git-town/git-town/pull/123 |
    And the current branch is "feature"
    And tool "open" is installed

  Scenario: a PR for this branch exists already
    When I run "git-town propose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                        |
      | feature | git fetch --prune --tags                                                       |
      |         | Finding proposal from feature into main ... #1 (Proposal from feature to main) |
      |         | open https://github.com/git-town/git-town/pull/123                             |
    And the initial branches and lineage exist now
    And the initial proposals exist now
  #
  # NOTE: Cannot test undo here.
  # The Git Town command under test has not created an undoable runstate.
  # Executing "git town undo" would undo the Git Town command executed during setup.
