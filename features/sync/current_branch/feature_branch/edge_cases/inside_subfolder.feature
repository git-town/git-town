Feature: sync inside a folder that doesn't exist on the main branch

  Background:
    Given my repo has the feature branches "alpha" and "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE       | FILE NAME        |
      | main   | local, remote | main commit   | main_file        |
      | alpha  | local, remote | folder commit | new_folder/file1 |
      | beta   | local, remote | beta commit   | file2            |
    And I am on the "alpha" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all" in the "new_folder" folder

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                          |
      | alpha  | git fetch --prune --tags         |
      |        | git add -A                       |
      |        | git stash                        |
      |        | git checkout main                |
      | main   | git rebase origin/main           |
      |        | git checkout alpha               |
      | alpha  | git merge --no-edit origin/alpha |
      |        | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout beta                |
      | beta   | git merge --no-edit origin/beta  |
      |        | git merge --no-edit main         |
      |        | git push                         |
      |        | git checkout alpha               |
      | alpha  | git push --tags                  |
      |        | git stash pop                    |
    And all branches are now synchronized
    And I am still on the "alpha" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE                        |
      | main   | local, remote | main commit                    |
      | alpha  | local, remote | folder commit                  |
      |        |               | main commit                    |
      |        |               | Merge branch 'main' into alpha |
      | beta   | local, remote | beta commit                    |
      |        |               | main commit                    |
      |        |               | Merge branch 'main' into beta  |
