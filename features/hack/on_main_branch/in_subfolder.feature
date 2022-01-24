Feature: git town-hack: starting a new feature from a subfolder on the main branch (with remote repo)

  As a developer working in a subfolder on the main branch
  I want to be able to extract my open changes into a feature branch
  So that I can get them reviewed.

  Background:
    Given the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE       | FILE NAME        |
      | main   | local    | folder commit | new_folder/file1 |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                      |
      | main        | git fetch --prune --tags     |
      |             | git add -A                   |
      |             | git stash                    |
      |             | git rebase origin/main       |
      |             | git push                     |
      |             | git branch new-feature main  |
      |             | git checkout new-feature     |
      | new-feature | git stash pop                |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE       |
      | main        | local, remote | folder commit |
      | new-feature | local         | folder commit |
