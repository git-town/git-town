Feature: syncing inside a folder that doesn't exist on the main branch

  Background:
    Given my repo has the feature branches "current-feature" and "other-feature"
    And my repo contains the commits
      | BRANCH          | LOCATION      | MESSAGE              | FILE NAME        |
      | main            | local, remote | main commit          | main_file        |
      | current-feature | local, remote | folder commit        | new_folder/file1 |
      | other-feature   | local, remote | other feature commit | file2            |
    And I am on the "current-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH          | COMMAND                                    |
      | current-feature | git fetch --prune --tags                   |
      |                 | git add -A                                 |
      |                 | git stash                                  |
      |                 | git checkout main                          |
      | main            | git rebase origin/main                     |
      |                 | git checkout current-feature               |
      | current-feature | git merge --no-edit origin/current-feature |
      |                 | git merge --no-edit main                   |
      |                 | git push                                   |
      |                 | git checkout other-feature                 |
      | other-feature   | git merge --no-edit origin/other-feature   |
      |                 | git merge --no-edit main                   |
      |                 | git push                                   |
      |                 | git checkout current-feature               |
      | current-feature | git push --tags                            |
      |                 | git stash pop                              |
    And I am still on the "current-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH          | LOCATION      | MESSAGE                                  |
      | main            | local, remote | main commit                              |
      | current-feature | local, remote | folder commit                            |
      |                 |               | main commit                              |
      |                 |               | Merge branch 'main' into current-feature |
      | other-feature   | local, remote | other feature commit                     |
      |                 |               | main commit                              |
      |                 |               | Merge branch 'main' into other-feature   |
