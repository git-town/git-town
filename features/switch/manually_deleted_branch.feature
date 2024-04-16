Feature: switch branches while a manually deleted branch is still listed in the lineage

  Scenario: repo contains a manually deleted branch
    Given the current branch is a local feature branch "alpha"
    And a local feature branch "beta"
    And a local feature branch "gamma"
    And I run "git branch -D beta"
    When I run "git-town switch" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it runs the commands
      | BRANCH | COMMAND            |
      | alpha  | git checkout gamma |
    And the current branch is now "gamma"
