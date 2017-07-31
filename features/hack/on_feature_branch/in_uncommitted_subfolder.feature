Feature: git town-hack: starting a new feature from a new subfolder on the main branch

  As a developer working on a feature branch that contains a subfolder that doesn't exist on the main branch
  I want to be able to extract my open changes into a new feature branch
  So that I can get them reviewed separately from the changes on this branch.


  Background:
    Given I have a feature branch named "feature"
    And the following commit exists in my repository
      | BRANCH | LOCATION         | MESSAGE     | FILE NAME |
      | main   | local and remote | main commit | main_file |
    And I am on the "feature" branch
    And I have an uncommitted file with name: "new_folder/file1" and content: "foo"
    When I run `git-town hack new-feature` in the "new_folder" folder


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                              |
      | feature     | git fetch --prune                    |
      | <none>      | cd <%= git_root_folder %>            |
      | feature     | git add -A                           |
      |             | git stash                            |
      |             | git checkout main                    |
      | main        | git rebase origin/main               |
      |             | git checkout -b new-feature main     |
      | new-feature | git stash pop                        |
      | <none>      | cd <%= git_root_folder %>/new_folder |
    And I end up on the "new-feature" branch
    And I am in the project root folder
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH      | LOCATION         | MESSAGE     |
      | main        | local and remote | main commit |
      | new-feature | local            | main commit |
