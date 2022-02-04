@skipWindows
Feature: sync before creating the pull request

  Background:
    Given my repo has a feature branch "parent"
    And my repo has a feature branch "child" as a child of "parent"
    And my repo contains the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | remote   | remote main commit   |
      | parent | local    | local parent commit  |
      |        | remote   | remote parent commit |
      | child  | local    | local child commit   |
      |        | remote   | remote child commit  |
    And my computer has the "open" tool installed
    And my repo's origin is "git@github.com:git-town/git-town.git"
    And I am on the "child" branch
    And my workspace has an uncommitted file
    When I run "git-town new-pull-request"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                                                   |
      | child  | git fetch --prune --tags                                                  |
      |        | git add -A                                                                |
      |        | git stash                                                                 |
      |        | git checkout main                                                         |
      | main   | git rebase origin/main                                                    |
      |        | git push                                                                  |
      |        | git checkout parent                                                       |
      | parent | git merge --no-edit origin/parent                                         |
      |        | git merge --no-edit main                                                  |
      |        | git push                                                                  |
      |        | git checkout child                                                        |
      | child  | git merge --no-edit origin/child                                          |
      |        | git merge --no-edit parent                                                |
      |        | git push                                                                  |
      |        | git stash pop                                                             |
      | <none> | open https://github.com/git-town/git-town/compare/parent...child?expand=1 |
    And "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent...child?expand=1
      """
    And I am still on the "child" branch
    And my workspace still contains my uncommitted file
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local, remote | remote main commit                                       |
      |        |               | local main commit                                        |
      | child  | local, remote | local child commit                                       |
      |        |               | remote child commit                                      |
      |        |               | Merge remote-tracking branch 'origin/child' into child   |
      |        |               | local parent commit                                      |
      |        |               | remote parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | remote main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
      |        |               | Merge branch 'parent' into child                         |
      | parent | local, remote | local parent commit                                      |
      |        |               | remote parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      |        |               | remote main commit                                       |
      |        |               | local main commit                                        |
      |        |               | Merge branch 'main' into parent                          |
