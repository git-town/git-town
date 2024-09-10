Feature: switch to branches described by a regex

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | alpha   | feature | main   | local     |
      | aloha   | feature | main   | local     |
      | another | feature | main   | local     |
      | beta    | feature | main   | local     |
    And the current branch is "aloha"
    When I run "git-town switch ^al" and enter into the dialogs:
      | KEYS       |
      | down enter |
    Then it runs the commands
      | BRANCH | COMMAND            |
      | aloha  | git checkout alpha |
    And the current branch is now "alpha"
