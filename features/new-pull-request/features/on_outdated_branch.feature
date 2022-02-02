Feature: syncing before creating the pull request

  Background:
    Given my code base has a feature branch "parent-feature"
    And my code base has a feature branch "child-feature" as a child of "parent-feature"
    And my repo contains the commits
      | BRANCH         | LOCATION | MESSAGE              |
      | main           | local    | local main commit    |
      |                | remote   | remote main commit   |
      | parent-feature | local    | local parent commit  |
      |                | remote   | remote parent commit |
      | child-feature  | local    | local child commit   |
      |                | remote   | remote child commit  |
    And my computer has the "open" tool installed
    And my repo's origin is "git@github.com:git-town/git-town.git"
    And I am on the "child-feature" branch
    And my workspace has an uncommitted file
    When I run "git-town new-pull-request"

  @skipWindows
  Scenario: result
    Then it runs the commands
      | BRANCH         | COMMAND                                                                                   |
      | child-feature  | git fetch --prune --tags                                                                  |
      |                | git add -A                                                                                |
      |                | git stash                                                                                 |
      |                | git checkout main                                                                         |
      | main           | git rebase origin/main                                                                    |
      |                | git push                                                                                  |
      |                | git checkout parent-feature                                                               |
      | parent-feature | git merge --no-edit origin/parent-feature                                                 |
      |                | git merge --no-edit main                                                                  |
      |                | git push                                                                                  |
      |                | git checkout child-feature                                                                |
      | child-feature  | git merge --no-edit origin/child-feature                                                  |
      |                | git merge --no-edit parent-feature                                                        |
      |                | git push                                                                                  |
      |                | git stash pop                                                                             |
      | <none>         | open https://github.com/git-town/git-town/compare/parent-feature...child-feature?expand=1 |
    Then "open" launches a new pull request with this url in my browser:
      """
      https://github.com/git-town/git-town/compare/parent-feature...child-feature?expand=1
      """
    And I am still on the "child-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH         | LOCATION      | MESSAGE                                                                  |
      | main           | local, remote | remote main commit                                                       |
      |                |               | local main commit                                                        |
      | child-feature  | local, remote | local child commit                                                       |
      |                |               | remote child commit                                                      |
      |                |               | Merge remote-tracking branch 'origin/child-feature' into child-feature   |
      |                |               | local parent commit                                                      |
      |                |               | remote parent commit                                                     |
      |                |               | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |
      |                |               | remote main commit                                                       |
      |                |               | local main commit                                                        |
      |                |               | Merge branch 'main' into parent-feature                                  |
      |                |               | Merge branch 'parent-feature' into child-feature                         |
      | parent-feature | local, remote | local parent commit                                                      |
      |                |               | remote parent commit                                                     |
      |                |               | Merge remote-tracking branch 'origin/parent-feature' into parent-feature |
      |                |               | remote main commit                                                       |
      |                |               | local main commit                                                        |
      |                |               | Merge branch 'main' into parent-feature                                  |
