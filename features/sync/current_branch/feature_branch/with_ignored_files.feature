Feature: syncing with ignored files

  As a developer using an IDE that creates a temp folder I don't want to check in
  I want "git-town sync" to leave those ignored files alone
  So that my IDE settings are not deleted while developing.

  - all files ignored by Git survive a "git sync" process unchanged

  Scenario: running "git sync" with ignored files
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE   | FILE NAME  | FILE CONTENT |
      | feature | local, remote | my commit | .gitignore | ignored      |
    And I am on the "feature" branch
    And my workspace has an uncommitted file with name "somefile" and content "important"
    And my workspace has an uncommitted file with name "test/ignored/important" and content "very important"
    When I run "git-town sync"
    Then my workspace still contains the file "test/ignored/important" with content "very important"
