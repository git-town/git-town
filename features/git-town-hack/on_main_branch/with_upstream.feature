Feature: git-hack: on the main branch with a upstream remote

  Background:
    Given my repo has an upstream repo
    And the following commits exist in my repository
      | BRANCH | LOCATION | MESSAGE         |
      | main   | upstream | upstream commit |
    And I am on the "main" branch
    And I have an uncommitted file
    When I run `gt hack new-feature`


  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                          |
      | main        | git fetch --prune                |
      |             | git add -A                       |
      |             | git stash                        |
      |             | git rebase origin/main           |
      |             | git fetch upstream               |
      |             | git rebase upstream/main         |
      |             | git push                         |
      |             | git checkout -b new-feature main |
      | new-feature | git push -u origin new-feature   |
      |             | git stash pop                    |
    And I am still on the "new-feature" branch
    And I still have my uncommitted file
    And I have the following commits
      | BRANCH      | LOCATION                    | MESSAGE         |
      | main        | local, remote, and upstream | upstream commit |
      | new-feature | local and remote            | upstream commit |
