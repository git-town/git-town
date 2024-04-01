Feature: compress the commits on a feature branch verbosely

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
      |         |               | commit 3 | file_3    | content 3    |
    And an uncommitted file
    When I run "git-town compress --verbose" and enter "compressed commit" for the commit message

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                            |
      |         | git version                                        |
      |         | git config -lz --global                            |
      |         | git config -lz --local                             |
      |         | git rev-parse --show-toplevel                      |
      |         | git rev-parse --verify --abbrev-ref @{-1}          |
      |         | git stash list                                     |
      |         | git status --long --ignore-submodules              |
      |         | git branch -vva --sort=refname                     |
      |         | git remote                                         |
      | feature | git fetch --prune --tags                           |
      | <none>  | git branch -vva --sort=refname                     |
      |         | git cherry -v main feature                         |
      | feature | git add -A                                         |
      |         | git stash                                          |
      |         | git reset --soft main                              |
      |         | git commit -m "commit 1"                           |
      | <none>  | git rev-list --left-right feature...origin/feature |
      | feature | git push --force-with-lease --force-if-includes    |
      | <none>  | git stash list                                     |
      | feature | git stash pop                                      |
      | <none>  | git branch -vva --sort=refname                     |
      |         | git config -lz --global                            |
      |         | git config -lz --local                             |
      |         | git stash list                                     |
    And it prints:
      """
      Ran 24 shell commands
      """
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE  |
      | feature | local, origin | commit 1 |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH  | COMMAND                                            |
      |         | git version                                        |
      |         | git config -lz --global                            |
      |         | git config -lz --local                             |
      |         | git rev-parse --show-toplevel                      |
      |         | git stash list                                     |
      |         | git status --long --ignore-submodules              |
      |         | git branch -vva --sort=refname                     |
      |         | git rev-parse --verify --abbrev-ref @{-1}          |
      |         | git remote get-url origin                          |
      | feature | git add -A                                         |
      |         | git stash                                          |
      | <none>  | git rev-parse --short HEAD                         |
      | feature | git reset --hard {{ sha 'commit 3' }}              |
      | <none>  | git rev-list --left-right feature...origin/feature |
      | feature | git push --force-with-lease --force-if-includes    |
      | <none>  | git stash list                                     |
      | feature | git stash pop                                      |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
