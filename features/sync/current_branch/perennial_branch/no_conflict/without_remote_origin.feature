Feature: git-town sync: syncing the current perennial branch (without remote repo)

  As a developer syncing a perennial branch and without a remote repository
  I want to be able update my ongoing work to include the latest finished features from the rest of the team
  So that our collaboration remains effective.

  Background:
    Given my repo does not have a remote origin
    And my repo has the local perennial branches "production" and "qa"
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE      | FILE NAME  |
      | main   | local    | main commit  | main_file  |
      | qa     | local    | local commit | local_file |
    And I am on the "qa" branch
    And my workspace has an uncommitted file
    When I run "git-town sync"

  Scenario: no conflict
    Then it runs the commands
      | BRANCH | COMMAND       |
      | qa     | git add -A    |
      |        | git stash     |
      |        | git stash pop |
    And I am still on the "qa" branch
    And my workspace still contains my uncommitted file
    And my repo is left with my original commits
