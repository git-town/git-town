@messyoutput
Feature: switch to branches described by several regexes

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | aloha   | feature | main   | local     |
      | alpha   | feature | main   | local     |
      | another | feature | main   | local     |
      | beta    | feature | main   | local     |
    And the current branch is "alpha"
    When I run "git-town switch ^al main" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then Git Town runs the commands
      | BRANCH | COMMAND           |
      | alpha  | git checkout main |
    And the current branch is now "main"
