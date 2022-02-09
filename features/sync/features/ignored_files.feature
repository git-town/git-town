Feature: ignore files

  Scenario: with ignored files
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE   | FILE NAME  | FILE CONTENT |
      | feature | local    | my commit | .gitignore | ignored      |
    And an uncommitted file with name "test/ignored/important" and content "changed ignored file"
    When I run "git-town sync"
    Then my workspace still contains the file "test/ignored/important" with content "changed ignored file"
