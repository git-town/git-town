Feature: display the local branch hierarchy in the middle of an ongoing rebase

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | alpha        | feature      | main   | local, origin |
      | beta         | feature      | alpha  | local, origin |
      | gamma        | feature      | beta   | local, origin |
      | conflicting  | feature      | main   | local, origin |
      | observed     | observed     |        | local, origin |
      | contribution | contribution |        | local, origin |
      | prototype    | prototype    | main   | local         |
      | parked       | parked       | main   | local         |
      | perennial    | perennial    |        | local, origin |
    And the commits
      | BRANCH      | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | conflicting | local    | conflicting local commit  | conflicting_file | local content  |
      |             | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "conflicting"
    And I ran "git pull --rebase"
    When I run "git-town branch"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
        main
          alpha
            beta
              gamma
      *   conflicting
          parked  (parked)
          prototype  (prototype)
        contribution  (contribution)
        observed  (observed)
        perennial  (perennial)
      """
