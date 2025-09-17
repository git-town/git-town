Feature: change offline mode at detached head

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | branch | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  |
      | branch | local, origin | commit 1 |
      |        | local, origin | commit 2 |
    And the current branch is "branch"
    And I ran "git checkout HEAD^"
    When I run "git-town offline 1"

  Scenario: result
    Then global Git setting "git-town.offline" is now "true"
