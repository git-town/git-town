Feature: git ship: aborting the ship of the supplied feature branch by entering an empty commit message

  (see ../../current_branch/on_feature_branch/without_open_changes/empty_commit_message.feature)


  Background:
    Given I have feature branches named "feature" and "other-feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | main    | local and remote | main commit    | main_file    | main content    |
      | feature | local            | feature commit | feature_file | feature content |
    And I am on the "other-feature" branch
    And I have an uncommitted file with name: "feature_file" and content: "conflicting content"
    When I run `git ship feature` and enter an empty commit message


  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                                      |
      | other-feature | git fetch --prune                            |
      |               | git stash -a                                 |
      |               | git checkout main                            |
      | main          | git rebase origin/main                       |
      |               | git checkout feature                         |
      | feature       | git merge --no-edit origin/feature           |
      |               | git merge --no-edit main                     |
      |               | git checkout main                            |
      | main          | git merge --squash feature                   |
      |               | git commit                                   |
      |               | git reset --hard                             |
      |               | git checkout feature                         |
      | feature       | git reset --hard <%= sha 'feature commit' %> |
      |               | git checkout main                            |
      | main          | git checkout other-feature                   |
      | other-feature | git stash pop                                |
    And I get the error "Ship aborted because commit exited with error"
    And I am still on the "other-feature" branch
    And I still have my uncommitted file
    And I am left with my original commits
