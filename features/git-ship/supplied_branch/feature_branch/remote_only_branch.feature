Feature: git ship: shipping the supplied feature branch with a tracking branch

  As a developer having finished a feature on another machine
  I want to be able to ship it without explicity fetching
  So that I can quickly move on to the next feature and remain productive.

  Background:
    Given I have a feature branch named "other-feature"
    And I have a feature branch named "feature" on another machine
    And the following commit exists in my repository on another machine
      | BRANCH  | LOCATION         | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local and remote | feature commit | feature_file | feature content |
    And I am on the "other-feature" branch
    And I have an uncommitted file with name: "feature_file" and content: "conflicting content"
    When I run `git ship feature -m "feature done"`


  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                            |
      | other-feature | git fetch --prune                  |
      |               | git stash -u                       |
      |               | git checkout main                  |
      | main          | git rebase origin/main             |
      |               | git checkout feature               |
      | feature       | git merge --no-edit origin/feature |
      |               | git merge --no-edit main           |
      |               | git checkout main                  |
      | main          | git merge --squash feature         |
      |               | git commit -m "feature done"       |
      |               | git push                           |
      |               | git push origin :feature           |
      |               | git branch -D feature              |
      |               | git checkout other-feature         |
      | other-feature | git stash pop                      |
    And I end up on the "other-feature" branch
    And I still have my uncommitted file
    And there is no "feature" branch
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE      | FILE NAME    |
      | main   | local and remote | feature done | feature_file |
