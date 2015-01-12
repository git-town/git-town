Feature: git ship: shipping the current feature branch

  As a developer entering a commit message that contains a double quote
  I want it to still work as expected
  So shipping is a robust process.

  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH  | LOCATION | FILE NAME    | FILE CONTENT    |
      | feature | remote   | feature_file | feature content |
    And I am on the "feature" branch
    When I run `git ship -m 'feature done with "double quotes"'`


  Scenario: result
    Then it runs the Git commands
      | BRANCH  | COMMAND                                             |
      | feature | git checkout main                                   |
      | main    | git fetch --prune                                   |
      | main    | git rebase origin/main                              |
      | main    | git checkout feature                                |
      | feature | git merge --no-edit origin/feature                  |
      | feature | git merge --no-edit main                            |
      | feature | git checkout main                                   |
      | main    | git merge --squash feature                          |
      | main    | git commit -m "feature done with \"double quotes\"" |
      | main    | git push                                            |
      | main    | git push origin :feature                            |
      | main    | git branch -D feature                               |
    And I end up on the "main" branch
    And there are no more feature branches
    And there are no open changes
    And I have the following commits
      | BRANCH | LOCATION         | MESSAGE                           | FILE NAME    |
      | main   | local and remote | feature done with "double quotes" | feature_file |
    And now I have the following committed files
      | BRANCH | FILES        |
      | main   | feature_file |
