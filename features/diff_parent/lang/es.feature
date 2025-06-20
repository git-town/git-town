Feature: view changes made on the current feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town diff-parent" with these environment variables
      | LANG | es_ES.UTF-8 |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git diff --merge-base main feature |
