Feature: ignoring files

  Scenario: running "git sync" with ignored files
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE   | FILE NAME  | FILE CONTENT |
      | feature | local    | my commit | .gitignore | ignored      |
    And I am on the "feature" branch
    And my workspace has an uncommitted file with name "test/ignored/important" and content "changed ignored file"
    When I run "git-town sync"
    Then my workspace still contains the file "test/ignored/important" with content "changed ignored file"
