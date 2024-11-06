Feature: display the local branch hierarchy with types for unusual branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | alpha        | feature      | main   | local, origin |
      | beta         | feature      | alpha  | local, origin |
      | gamma        | feature      | beta   | local, origin |
      | observed     | observed     |        | local, origin |
      | contribution | contribution |        | local, origin |
      | prototype    | prototype    | main   | local         |
      | parked       | parked       | main   | local         |
      | perennial    | perennial    |        | local, origin |
    And the current branch is "beta"
    When I run "git-town branch"

  Scenario: result
    Then Git Town runs no commands
    And it prints:
      """
        main
          alpha
      *     beta
              gamma
          parked  (parked)
          prototype  (prototype)
        contribution  (contribution)
        observed  (observed)
        perennial  (perennial)
      """
