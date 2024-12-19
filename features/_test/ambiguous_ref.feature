@this
Feature: The tests succeed when dealing with an ambiguous ref name

  Scenario:
    Given a Git repo with origin
    And the branches
      | NAME        | TYPE    | PARENT | LOCATIONS |
      | <ambiguous> | feature | main   | local     |
    And the commits
      | BRANCH      | LOCATION | MESSAGE       |
      | <ambiguous> | local    | spooky commit |
    When I run "git-town"
    Then these commits exist now
      | BRANCH      | LOCATION | MESSAGE       |
      | <ambiguous> | local    | spooky commit |
