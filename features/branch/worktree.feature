Feature: display the local branch hierarchy in the middle of an ongoing rebase

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
    And branch "gamma" is active in another worktree
    And I ran "git pull --rebase"
    When I run "git-town branch"

  Scenario: result
    Then it runs no commands
    And it prints:
      """
        main
          alpha
      *     beta
      +       gamma
          parked  (parked)
          prototype  (prototype)
        contribution  (contribution)
        observed  (observed)
        perennial  (perennial)
      """
