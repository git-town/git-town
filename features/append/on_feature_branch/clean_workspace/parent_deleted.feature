Feature: append a branch to a branch whose parent was shipped on the remote

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent |               |
    And the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And origin ships the "parent" branch
    And inspect the repo
    And the current branch is "child"
    When I run "git-town append new -v"

  @debug @this
  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                        |
      |        | git version                                    |
      |        | git rev-parse --show-toplevel                  |
      |        | git config -lz --includes --global             |
      |        | git config -lz --includes --local              |
      |        | git status --long --ignore-submodules          |
      |        | git remote                                     |
      |        | git rev-parse --abbrev-ref HEAD                |
      | child  | git fetch --prune --tags                       |
      | <none> | git stash list                                 |
      |        | git branch -vva --sort=refname                 |
      |        | git rev-parse --verify --abbrev-ref @{-1}      |
      | child  | git checkout main                              |
      | main   | git rebase origin/main                         |
      | <none> | git rev-list --left-right main...origin/main   |
      | main   | git checkout parent                            |
      | parent | git merge --no-edit --ff main                  |
      | <none> | git diff main..parent                          |
      | parent | git checkout child                             |
      | child  | git merge --no-edit --ff origin/child          |
      |        | git merge --no-edit --ff parent                |
      | <none> | git rev-list --left-right child...origin/child |
      | child  | git push                                       |
      | <none> | git show-ref --verify --quiet refs/heads/child |
      | child  | git checkout -b new                            |
      | <none> | git show-ref --verify --quiet refs/heads/child |
      |        | git config git-town-branch.new.parent child    |
      |        | git show-ref --verify --quiet refs/heads/child |
      |        | git branch -vva --sort=refname                 |
      |        | git config -lz --includes --global             |
      |        | git config -lz --includes --local              |
      |        | git stash list                                 |
    And it prints:
      """
      deleted branch "parent"
      """
    And the current branch is now "new"
    And the branches are now
      | REPOSITORY | BRANCHES         |
      | local      | main, child, new |
      | origin     | main, child      |
    And this lineage exists now
      | BRANCH | PARENT |
      | child  | main   |
      | new    | child  |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | new    | git checkout child                              |
      | child  | git reset --hard {{ sha 'child commit' }}       |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout main                               |
      | main   | git reset --hard {{ sha 'initial commit' }}     |
      |        | git branch parent {{ sha 'parent commit' }}     |
      |        | git checkout child                              |
      | child  | git branch -D new                               |
    And the current branch is still "child"
    And the initial branches and lineage exist
