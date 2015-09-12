Feature: git extract: extracting multiple commits (without remote repo)

  (see ../with_remote_origin/with_open_changes.feature)


  Background:
    Given I have a feature branch named "feature"
    And my repo does not have a remote origin
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE          |
      | main    | local    | main commit      |
      | feature | local    | feature commit   |
      |         |          | refactor1 commit |
      |         |          | refactor2 commit |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git extract refactor` with the last two commit shas


  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                                                                     |
      | feature  | git stash -u                                                                |
      |          | git checkout main                                                           |
      | main     | git checkout -b refactor main                                               |
      | refactor | git cherry-pick <%= sha 'refactor1 commit' %> <%= sha 'refactor2 commit' %> |
      |          | git stash pop                                                               |
    And I end up on the "refactor" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH   | LOCATION | MESSAGE          |
      | main     | local    | main commit      |
      | feature  | local    | feature commit   |
      |          |          | refactor1 commit |
      |          |          | refactor2 commit |
      | refactor | local    | main commit      |
      |          |          | refactor1 commit |
      |          |          | refactor2 commit |
