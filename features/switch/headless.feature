@messyoutput
Feature: switch branches from detached head

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  |
      | alpha  | local, origin | commit 1 |
      |        | local, origin | commit 2 |
    And the current branch is "alpha"
    And I ran "git checkout HEAD^"
    When I run "git-town switch" and enter into the dialogs:
      | DIALOG        | KEYS     |
      | switch-branch | up enter |

  Scenario: result
    Then Git Town runs the commands
      | BRANCH               | COMMAND            |
      | {{ sha 'commit 1' }} | git checkout alpha |
