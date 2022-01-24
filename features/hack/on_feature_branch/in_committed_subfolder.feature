Feature: git town-hack: starting a new feature from a new subfolder on the main branch

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE       | FILE NAME        |
      | main    | local, remote | main commit   | main_file        |
      | feature | local, remote | folder commit | new_folder/file1 |
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature" in the "new_folder" folder


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | feature     | git fetch --prune --tags    |
      |             | git add -A                  |
      |             | git stash                   |
      |             | git checkout main           |
      | main        | git rebase origin/main      |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE       |
      | main        | local, remote | main commit   |
      | feature     | local, remote | folder commit |
      | new-feature | local         | main commit   |
