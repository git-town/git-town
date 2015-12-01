Feature: git hack: starting a new feature from a new subfolder on the main branch

  As a developer working on a feature branch that contains a subfolder that doesn't exist on the main branch
  I want to be able to extract my open changes into a new feature branch
  So that I can get them reviewed separately from the changes on this branch.


  This feature cannot be tested.
  GT performs all the correct commands.
  But everything happens in a subshell.
  Git removes the folder that the user session is currently in,
  then later creates a new folder with the same name.
  The user session doesn't know that, so it is now in a folder that doesn't exist.

  When encountering this issue, simply cd into the Git root folder after git-hack is done.

  Strangely enough, the same test (with the same side effects) works for git sync:
  https://github.com/Originate/git-town/blob/master/features/git-sync/folder_does_not_exist_on_main_branch/no_conflict.feature
  Maybe creating a completely new Git branch is what breaks things here.


  Background:
    Given I have a feature branch named "feature"
    Given the following commit exists in my repository
      | BRANCH  | LOCATION         | MESSAGE       | FILE NAME        |
      | main    | local and remote | main commit   | main_file        |
      | feature | local and remote | folder commit | new_folder/file1 |
    And I am on the "feature" branch
    And I have an uncommitted file
    When I run `git hack new-feature` in the "new_folder" folder


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                           |
      | feature     | git fetch --prune                 |
      | <none>      | cd <%= git_root_folder %>         |
      | feature     | git stash -u                      |
      |             | git checkout main                 |
      | main        | git rebase origin/main            |
      |             | git checkout -b new-feature main  |
      | new-feature | git push -u origin new-feature    |
      |             | git stash pop                     |
    And I end up on the "new-feature" branch
    And I am in the project root folder
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH      | LOCATION         | MESSAGE       |
      | main        | local and remote | main commit   |
      | feature     | local and remote | folder commit |
      | new-feature | local and remote | main commit   |
