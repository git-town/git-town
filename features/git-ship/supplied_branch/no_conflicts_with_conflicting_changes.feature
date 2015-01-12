Feature: git ship: shipping the supplied feature branch (with conflicting changes)

  As a developer getting the LGTM for a feature branch while working on unrelated things that conflict with the main branch
  I want to be able to ship the approved branch anyways
  So that I don't have to execute a bunch of boilerplate Git commands to ship, and remain productive and focussed on my current work.


  Scenario: local feature branch
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | main    | local    | main commit    | main_file    | main content    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "main_file" and content: "conflicting content"
    When I run `git ship feature -m 'feature done'`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git stash -u                       |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
      | main          | git push                           |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      | feature       | git merge --no-edit main           |
      | feature       | git checkout main                  |
      | main          | git merge --squash feature         |
      | main          | git commit -m 'feature done'       |
      | main          | git push                           |
      | main          | git push origin :feature           |
      | main          | git branch -D feature              |
      | main          | git checkout other_feature         |
      | other_feature | git stash pop                      |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "main_file" and content: "conflicting content"
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | main commit  | main_file    |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES                   |
      | main   | feature_file, main_file |


  Scenario: feature branch with non-pulled updates in the repo
    Given I have feature branches named "feature" and "other_feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION         | FILE NAME    | FILE CONTENT          |
      | feature | local and remote | feature_file | early feature content |
      | feature | local and remote | feature_file | mid feature content   |
      | feature | remote           | feature_file | final feature content |
    And I am on the "other_feature" branch
    And I have an uncommitted file with name: "feature_file" and content: "conflicting content"
    When I run `git ship feature -m 'feature done'`
    Then it runs the Git commands
      | BRANCH        | COMMAND                            |
      | other_feature | git stash -u                       |
      | other_feature | git checkout main                  |
      | main          | git fetch --prune                  |
      | main          | git rebase origin/main             |
      | main          | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      | feature       | git merge --no-edit main           |
      | feature       | git checkout main                  |
      | main          | git merge --squash feature         |
      | main          | git commit -m 'feature done'       |
      | main          | git push                           |
      | main          | git push origin :feature           |
      | main          | git branch -D feature              |
      | main          | git checkout other_feature         |
      | other_feature | git stash pop                      |
    And I end up on the "other_feature" branch
    And I still have an uncommitted file with name: "feature_file" and content: "conflicting content"
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
