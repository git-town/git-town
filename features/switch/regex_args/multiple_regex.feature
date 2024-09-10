Feature: switch to branches described by several regexes

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | alpha   | feature | main   | local     |
      | aloha   | feature | main   | local     |
      | another | feature | main   | local     |
      | beta    | feature | main   | local     |
    And the current branch is "aloha"
    When I run "git-town switch ^al main" and enter into the dialogs:
      | KEYS     |
      | up enter |
    Then it runs the commands
      | BRANCH | COMMAND           |
      | aloha  | git checkout main |
    And the current branch is now "main"
